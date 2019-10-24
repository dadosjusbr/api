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

var netClient = &http.Client{
	Timeout: time.Second * 60,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func main() {
	month := flag.Int("mes", 0, "Mês a ser analisado")
	year := flag.Int("ano", 0, "Ano a ser analisado")
	name := flag.String("nome", "", "Used for login purposes")
	cpf := flag.String("cpf", "", "used for login purpose. format xxx.xxx.xxx-xx")
	flag.Parse()
	if *month == 0 || *year == 0 || *cpf == "" || *name == "" {
		log.Fatalf("Need all arguments to continue, please try again\n")
	}

	acessCode, err := login(*name, *cpf)
	if err != nil {
		log.Fatalf("login error: %q", err)
	}

	data, err := queryData(acessCode, *month, *year)
	if err != nil {
		log.Fatalf("Query data error: %q", err)
	}

	dataDesc := fmt.Sprintf("remuneracoes-trepb-%02d-%04d", *month, *year)

	if err = save(dataDesc, data); err != nil {
		log.Fatalf("Error saving data to file: %q", err)
	}

}

// save downloads content from url and save it on a file.
func save(desc string, data []*html.Node) error {
	fileName := fmt.Sprintf("%s.html", desc)
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file(%s):%q", fileName, err)
	}
	defer f.Close()

	for _, node := range data {
		nodeReader := strings.NewReader(htmlquery.OutputHTML(node, true))
		if io.Copy(f, nodeReader); err != nil {
			os.Remove(fileName)
			return fmt.Errorf("error copying response content to file: %q", err)
		}
	}
	return nil
}

// queryData query server for data of a specified month and year.
func queryData(acessCode string, month, year int) ([]*html.Node, error) {
	query := fmt.Sprintf(`acao=AnexoVIII&folha=&valida=true&toExcel=false&chaveDeAcesso=%s&mes=%d&ano=%04d`, acessCode, month, year)
	queryURL := fmt.Sprintf(`http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?%s`, query)
	doc, err := loadURL(queryURL)
	if err != nil {
		return nil, fmt.Errorf("error while loading url: %q", err)
	}

	tables, err := htmlquery.QueryAll(doc, "//table")
	if err != nil {
		return nil, fmt.Errorf("error while making query for data tables: %q", err)
	}
	if len(tables) == 0 {
		return nil, fmt.Errorf("couldn't find any data tables")
	}
	return tables, nil
}

//loadURL loads HTML document from specified URL.
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
