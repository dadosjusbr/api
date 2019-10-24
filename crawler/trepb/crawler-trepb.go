package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	baseURL        = "http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?acao=Anexo_VIII"
	questionXpath  = "/html/body/form/table/tbody/tr[2]/td/table/tbody/tr[3]/td/table/tbody/tr[3]/td[1]"
	acessCodeXpath = "/html/body/form/input[5]"
)

var monthStr = []string{"janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"}

var netClient = &http.Client{
	Timeout: time.Second * 60,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	flag.Parse()
	if *month == 0 || *year == 0 {
		log.Fatalf("need arguments to continue, please try again")
	}

	acessCode, err := login()
	if err != nil {
		log.Fatalf("login error: %q", err)
	}
	fmt.Println(acessCode)

}

//Load HTML document from specified URL.
func loadURL(baseURL string) (*html.Node, error) {
	resp, err := netClient.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %q", baseURL, err)
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading doc (%s): %q", baseURL, err)
	}
	return doc, nil
}

// Returns accessCode for the api.
func login() (string, error) {
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

	resp, err := loginRequest(question, ans)
	if err != nil {
		return "", fmt.Errorf("error while trying to make a login request: %q", err)
	}

	code := retrieveAcessCode(resp)
	if code == "" {
		return "", fmt.Errorf("couldn't retrieve access code")
	}

	return code, nil
}

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

func loginRequest(question, ans string) (io.Reader, error) {
	body := fmt.Sprintf(
		`nomeUsuario=Marcos+Barros+de+Medeiros+Filho&cpfUsuario=097.650.704-89&respostaCaptcha=%s&btnLogin=Efetuar+login&identificaUsuario=&perguntaCaptcha=%s`,
		url.QueryEscape(ans), url.QueryEscape(question))

	req, err := http.NewRequest("GET", fmt.Sprintf("http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?%s", body), nil)
	if err != nil {
		return nil, fmt.Errorf("error while trying to make a NewRequest structure: %q", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:69.0) Gecko/20100101 Firefox/69.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?acao=Anexo_VIII")
	req.Header.Set("Cookie", "JSESSIONID=197709FD583A7E6145E01453E36CAED9")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := netClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while trying to make the post request: %q", err)
	}
	defer resp.Body.Close()

	return resp.Body, nil
}

func saveDebbug(body io.Reader) {
	out, err := os.Create("filename5.html")
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, body)
}

// retrieveAcessCode searchs for accessCode inside the page and return if found.
func retrieveAcessCode(page io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(page)
	pageStr := buf.String()

	code := substringBetween(pageStr, `<input type="hidden" name="chaveDeAcesso" value="`, `"`)

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
