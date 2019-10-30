package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	baseURL        = "http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?acao=Anexo_VIII"
	questionXpath  = "/html/body/form/table/tbody/tr[2]/td/table/tbody/tr[3]/td/table/tbody/tr[3]/td[1]"
	acessCodeXpath = "/html/body/form/input[5]"
	acessCodeCache = "acessCode.txt"
)

// login returns the accessCode for the api.
func login(name, cpf string) (string, error) {
	doc, err := loadURL(baseURL)
	if err != nil {
		return "", fmt.Errorf("Error while trying to load document: %q", err)
	}

	question, err := findQuestion(doc)
	if err != nil {
		return "", fmt.Errorf("Error while trying to retrieve question from page: %q", err)
	}

	ans, err := solution(question)
	if err != nil {
		return "", fmt.Errorf("Error while trying to find answer to question: %q", err)
	}

	resp, err := loginRequest(question, ans, name, cpf)
	if err != nil {
		return "", fmt.Errorf("error while trying to make a login request: %q", err)
	}

	code := retrieveAcessCode(resp)
	if code == "" {
		return "", fmt.Errorf("couldn't retrieve access code. Question: %s. Answer: %s", question, ans)
	}

	return code, nil
}

func retrieveCachedCode() (string, error) {
	_, err := os.Stat(acessCodeCache)
	if err != nil {
		return "", err
	}

	f, err := os.Open(acessCodeCache)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return "", err
	}
	accessCode := buf.String()
	fmt.Println(accessCode)
	ok, err := validateKey(accessCode)
	if err != nil {
		return "", fmt.Errorf("error while validating key from cache file: %q", err)
	}
	if ok {
		return accessCode, nil
	}

	return "", nil
}

func validateKey(key string) (bool, error) {
	query := fmt.Sprintf(`acao=AnexoVIII&folha=&valida=true&toExcel=false&chaveDeAcesso=%s&mes=6&ano=2005`, key)
	requestURL := fmt.Sprintf("http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?%s", query)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return false, fmt.Errorf("error creating GET request to %s: %q", requestURL, err)
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making GET request to %s: %q", requestURL, err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return false, err
	}
	if err != nil {
		return false, fmt.Errorf("error while reading response body %s: %q", requestURL, err)
	}

	return len(buf.String()) != 0, nil
}

// findQuestion makes an xpath query to find captcha question inside the html page.
func findQuestion(doc *html.Node) (string, error) {
	qNode, err := htmlquery.Query(doc, questionXpath)
	if err != nil {
		return "", fmt.Errorf("Couldn't find Question: %q", err)
	}
	if qNode == nil {
		return "", fmt.Errorf("Couldn't find Question")
	}

	question := strings.TrimSpace(qNode.FirstChild.Data)
	return question, nil
}

//loginRequest makes login request and returns response body as a string
func loginRequest(question, ans, name, cpf string) (string, error) {
	body := fmt.Sprintf(
		`nomeUsuario=%s&cpfUsuario=%s&respostaCaptcha=%s&btnLogin=Efetuar+login&identificaUsuario=&perguntaCaptcha=%s`,
		url.QueryEscape(name), cpf, url.QueryEscape(ans), url.QueryEscape(url.QueryEscape(question)))

	req, err := http.NewRequest("GET", fmt.Sprintf("http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?%s", body), nil)
	if err != nil {
		return "", fmt.Errorf("error while trying to make a NewRequest structure: %q", err)
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?acao=Anexo_VIII")

	resp, err := netClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while trying to make the post request: %q", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	pageStr := buf.String()

	return pageStr, nil
}

// retrieveAcessCode searchs for accessCode inside the page and return if found.
func retrieveAcessCode(page string) string {
	code := substringBetween(page, `<input type="hidden" name="chaveDeAcesso" value="`, `"`)
	if len(code) != 32 {
		return ""
	}
	return code
}

//substringBetween returns the substring in str between before and after strings.
func substringBetween(str, before, after string) string {
	a := strings.SplitAfterN(str, before, 2)
	b := strings.SplitAfterN(a[len(a)-1], after, 2)
	if 1 == len(b) {
		return b[0]
	}
	return b[0][0 : len(b[0])-len(after)]
}
