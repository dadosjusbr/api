package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
		accessCode, err = newAccessCode(name, cpf)
		if err != nil {
			return "", fmt.Errorf("login error: %q", err)
		}
	}

	if err = validateKey(accessCode); err != nil {
		return "", fmt.Errorf("error while validating key from cache file: %q", err)
	}

	return accessCode, nil
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
		return "", fmt.Errorf("error openning cache file(%s): %q", cacheFileName, err)
	}
	defer f.Close()

	ac, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("error reading access code from cache file (%s): %q", cacheFileName, err)
	}
	return string(ac), nil
}

// newAccessCode makes an attempt to get a new access code from the api.
func newAccessCode(name, cpf string) (string, error) {
	question, err := captchaQuestion()
	if err != nil {
		return "", fmt.Errorf("error while trying to retrieve question: %q", err)
	}

	ans, err := solution(question)
	if err != nil {
		return "", fmt.Errorf("error while trying to find answer to question: %q", err)
	}

	code, err := login(question, ans, name, cpf)
	if err != nil {
		return "", fmt.Errorf("error trying to login: %q", err)
	}

	return code, nil
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

// captchaQuestion returns the captcha question.
func captchaQuestion() (string, error) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating get request to %s: %q", baseURL, err)
	}

	doc, err := httpReq(req)
	if err != nil {
		return "", fmt.Errorf("Error while trying to load document: %q", err)
	}

	question, err := findQuestion(doc, questionXpath)
	if err != nil {
		return "", fmt.Errorf("Error while trying to retrieve question from page: %q", err)
	}
	return question, nil
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

// login makes a login request and returns an accessCode if found.
func login(question, ans, name, cpf string) (string, error) {
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

// loginRequest makes login request and returns response body as a string
func loginRequest(question, ans, name, cpf string) (string, error) {
	body := fmt.Sprintf(
		`nomeUsuario=%s&cpfUsuario=%s&respostaCaptcha=%s&btnLogin=Efetuar+login&identificaUsuario=&perguntaCaptcha=%s`,
		url.QueryEscape(name), cpf, url.QueryEscape(ans), url.QueryEscape(url.QueryEscape(question)))
	reqURL := fmt.Sprintf("http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?%s", body)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("error while trying to make a NewRequest structure: %q", err)
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?acao=Anexo_VIII")

	doc, err := httpReq(req)
	if err != nil {
		return "", fmt.Errorf("Error while executing Get request to %s: %q", reqURL, err)
	}
	return htmlquery.OutputHTML(doc, true), nil
}

// saveToCache saves new code to cache file.
func saveToCache(code, cacheFileName string) error {
	f, err := os.Create(cacheFileName)
	if err != nil {
		return fmt.Errorf("error creating cache file: %q", err)
	}
	defer f.Close()
	n, err := f.Write([]byte(code))
	if err != nil {
		return fmt.Errorf("error writing to cache file: %q", err)
	}
	if n != len(code) {
		return fmt.Errorf("error writing code to cache file: Size of code is different from number of bytes written")
	}

	return nil
}

// retrieveAcessCode searchs for accessCode inside the page and return if found.
func retrieveAcessCode(page string) string {
	code := substringBetween(page, `<input type="hidden" name="chaveDeAcesso" value="`, `"`)
	if len(code) != 32 {
		return ""
	}
	return code
}
