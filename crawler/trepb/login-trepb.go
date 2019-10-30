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
	baseURL         = "http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?acao=Anexo_VIII"
	questionXpath   = "/html/body/form/table/tbody/tr[2]/td/table/tbody/tr[3]/td/table/tbody/tr[3]/td[1]"
	accessCodeXpath = "/html/body/form/input[5]"
	accessCodeCache = "acessCode.txt"
)

// accessCode gets access code from cache file or a new one from TRE-PB server.
func accessCode(name, cpf string) (string, error) {
	accessCode, err := retrieveCachedCode(accessCodeCache)
	if err != nil {
		return "", fmt.Errorf("retrieve cached code error: %q", err)
	}

	if accessCode == "" {
		accessCode, err = login(name, cpf)
		if err != nil {
			return "", fmt.Errorf("login error: %q", err)
		}
	}

	if err = validateKey(accessCode); err != nil {
		return "", fmt.Errorf("error while validating key from cache file: %q", err)
	}

	return accessCode, nil
}

// login returns the accessCode for the api.
func login(name, cpf string) (string, error) {
	doc, err := loadURL(baseURL)
	if err != nil {
		return "", fmt.Errorf("Error while trying to load document: %q", err)
	}

	question, err := findQuestion(doc, questionXpath)
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

	err = saveToCache(code, accessCodeCache)
	if err != nil {
		return "", fmt.Errorf("error while saving code to cache file: %q", err)
	}

	return code, nil
}

// retrieveChachedCode makes an attempt to retrieve a cached access code.
func retrieveCachedCode(cacheFileName string) (string, error) {
	if _, err := os.Stat(cacheFileName); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	f, err := os.Open(cacheFileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(f); err != nil {
		return "", err
	}
	accessCode := buf.String()

	return accessCode, nil
}

// validateKey makes a query to the TRE-PB API to assure key is valid.
func validateKey(key string) error {
	query := fmt.Sprintf(`acao=AnexoVIII&folha=&valida=true&toExcel=false&chaveDeAcesso=%s&mes=6&ano=2005`, key)
	requestURL := fmt.Sprintf("http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?%s", query)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("error creating GET request to %s: %q", requestURL, err)
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := netClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making GET request to %s: %q", requestURL, err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return fmt.Errorf("error while reading response body %s: %q", requestURL, err)
	}

	if len(buf.String()) == 0 {
		return fmt.Errorf("not a valid key: %s", key)
	}

	return nil
}

func saveToCache(code, cacheFileName string) error {
	f, err := os.Create(cacheFileName)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating cache file: %q", err)
	}
	defer f.Close()
	_, err = f.Write([]byte(code))
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error writing to cache file: %q", err)
	}
	return nil
}

// findQuestion makes an xpath query to find captcha question inside the html page.
func findQuestion(doc *html.Node, xpath string) (string, error) {
	qNode, err := htmlquery.Query(doc, xpath)
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
