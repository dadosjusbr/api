package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/antchfx/htmlquery"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

const htmlSample = `
<ul id="arquivos-2011">
	<li><a href="">Fevereiro 2011</a></li>
</ul>
<ul id="arquivos-2013-mes-01">
	<li><a href="servidores.pdf"></a></li>
	<li><a href="magistrados.pdf"></a></li>
</ul>
`

//Test if loadURL is loading the html doc.
func TestLoadURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, htmlSample)
	}))
	defer ts.Close()

	_, err := loadURL(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
}

// Test if xpath query is finding the interest nodes.
func TestFindInterestNodes(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(htmlSample))
	if err != nil {
		t.Fatal(err)
	}

	data := []struct {
		desc     string
		month    int
		year     int
		node     *html.Node
		nodeList []string
	}{
		{"Nodes past 2012", 1, 2013, doc, []string{
			`<a href="servidores.pdf"></a>`,
			`<a href="magistrados.pdf"></a>`,
		}},
		{"Nodes before 2013", 2, 2011, doc, []string{
			`<a href="">Fevereiro 2011</a>`,
		}},
	}

	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			got, err := findInterestNodes(d.node, d.month, d.year)
			assert.NoError(t, err)
			var nodeList []string
			for _, node := range got {
				nodeList = append(nodeList, htmlquery.OutputHTML(node, true))
			}
			assert.Equal(t, d.nodeList, nodeList)
		})
	}
}

// Test if interestNodes() returns an error if no node is found.
func TestFindInterestNodes_Error(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(htmlSample)) //LOADING DOC
	if err != nil {
		t.Fatal(err)
	}

	data := []struct {
		desc      string
		month     int
		year      int
		node      *html.Node
		errorDesc string
	}{
		{"nodes for given month and year not available", 1, 2015, doc, "couldn't find any link for 01-2015"},
	}

	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			_, err := findInterestNodes(d.node, d.month, d.year)
			assert.Error(t, err)
			assert.Equal(t, d.errorDesc, err.Error())
		})
	}
}

// Test if file name is returning appropriate names for the files.
func TestFileName(t *testing.T) {
	data := []struct {
		desc   string
		month  int
		year   int
		href   string
		result string
	}{
		{"Default name", 2, 2011, "anexo_viii_fev_20111.pdf", "remuneracoes-tjpb-02-2011"},
		{"Magistrados", 1, 2013, "magistrados.pdf", "remuneracoes-magistrados-tjpb-01-2013"},
		{"Servidor", 1, 2013, "servidores.pdf", "remuneracoes-servidores-tjpb-01-2013"},
	}

	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			assert.Equal(t, d.result, fileName(d.href, d.month, d.year))
		})
	}
}

// Test if the result of the request is saved in the buffer.
func TestDownload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	}))
	defer ts.Close()

	var buf bytes.Buffer
	assert.NoError(t, download(ts.URL, &buf))
	assert.Equal(t, "Hello", buf.String())
}

// Test if a file with the result is created. Download should asure content is the same.
func TestSave(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	}))
	defer ts.Close()

	assert.NoError(t, save("testFile", ts.URL))
	assert.FileExists(t, "testFile.pdf")
	assert.NoError(t, os.Remove("testFile.pdf"))
}

// Test if the file is erased if save returns an error.
func TestSave_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	err := save("testFile", ts.URL)
	assert.Error(t, err)
	_, err = os.Stat("testFile.pdf")
	assert.Error(t, err)
}
