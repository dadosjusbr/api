package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	baseURL           = "http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?acao=Anexo_VIII"
	captchaXpath      = "/html/body/form/table/tbody/tr[2]/td/table/tbody/tr[3]/td/table/tbody/tr[3]/td[1]"
	questionFormXpath = "//*[@name='perguntaCaptcha']"
)

var monthStr = []string{"janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"}

var netClient = &http.Client{
	Timeout: time.Second * 60,
}

func main() {
	q, err := login()
	if err != nil {
		log.Fatalf("%q", err)
	}
	fmt.Println(q)
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
	ans, err := solveCaptcha(doc)
	if err != nil {
		return "", fmt.Errorf("Error while trying to find answer to question: %q", err)
	}
	_, err = formatedQuestion(doc)
	if err != nil {
		return "", fmt.Errorf("Error while trying to find form question: %q", err)
	}

	return ans, nil
}

// Find question element that should be sent in the login request
func formatedQuestion(doc *html.Node) (string, error) {
	qFormNode, err := htmlquery.Query(doc, questionFormXpath)
	if err != nil {
		return "", fmt.Errorf("Couldn't find Form Question Node: %q", err)
	}
	if qFormNode == nil {
		return "", fmt.Errorf("Couldn't find Form Question Node")
	}

	return qFormNode.Attr[len(qFormNode.Attr)-1].Val, nil
}

// Find question and return answer.
func solveCaptcha(doc *html.Node) (string, error) {
	qNode, err := htmlquery.Query(doc, captchaXpath)
	if err != nil {
		return "", fmt.Errorf("Couldn't find Question: %q", err)
	}
	if qNode == nil {
		return "", fmt.Errorf("Couldn't find Question")
	}

	question := strings.TrimSpace(qNode.FirstChild.Data)
	ans, err := solution(question)
	if err != nil {
		return "", fmt.Errorf("Couldn't find solution for question: %q", err)
	}
	return ans, nil
}
