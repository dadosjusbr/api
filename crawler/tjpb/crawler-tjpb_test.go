package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const htmlSample = `<!DOCTYPE html><html lang="en-US">
<head>
</head>
<body>
<div>
<ul id="arquivos-2011" class="collapse in" aria-expanded="true" style="">
<li><a href="https://www.tjpb.jus.br/sites/default/files/anexos/2018/06/anexo_viii_fev_20111.pdf">Anexo VIII - Res. 102 CNJ - Fevereiro 2011</a></li>
</ul>
<ul id="arquivos-2013-mes-01" class="collapse">
<li><a href="https://www.tjpb.jus.br/sites/default/files/anexos/2018/06/201301_servidores.pdf">Anexo único - Res. 151 CNJ - Janeiro 2013 - Servidores</a></li>
<li><a href="https://www.tjpb.jus.br/sites/default/files/anexos/2018/06/201301_magistrados.pdf">Anexo único - Res. 151 CNJ - Janeiro 2013 - Magistrados</a></li>
</ul>
</div>
</body>
</html>
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

// Test if xpath query is finding the interest nodes for years past 2012
func TestInterestNodesPast2012(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, htmlSample)
	}))
	defer ts.Close()

	doc, _ := loadURL(ts.URL)
	nodes, err := interestNodes(doc, 1, 2013)
	assert.NoError(t, err)
	assert.NotEmpty(t, nodes)
	assert.Len(t, nodes, 2)
}

// Test if xpath query is finding the interest nodes for years before 2013
func TestInterestNodesBefore2013(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, htmlSample)
	}))
	defer ts.Close()

	doc, _ := loadURL(ts.URL)
	nodes, err := interestNodes(doc, 2, 2011)
	assert.NoError(t, err)
	assert.NotEmpty(t, nodes)
}

// Test if interestNodes() returns an error if no node is found.
func TestInterestNodesError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, htmlSample)
	}))
	defer ts.Close()

	doc, _ := loadURL(ts.URL)

	_, err := interestNodes(doc, 13, 2013)
	assert.Error(t, err)
}

// Test if file name is returning appropriate names for the files.
func TestFileName(t *testing.T) {
	def := fileName("https://www.tjpb.jus.br/sites/default/files/anexos/2018/06/anexo_viii_fev_20111.pdf", 2, 2011)
	assert.Equal(t, "remuneracoes-tjpb-02-2011", def)
	mag := fileName("https://www.tjpb.jus.br/sites/default/files/anexos/2018/06/201301_magistrados.pdf", 1, 2013)
	assert.Equal(t, "remuneracoes-magistrados-tjpb-01-2013", mag)
	serv := fileName("https://www.tjpb.jus.br/sites/default/files/anexos/2018/06/201301_servidores.pdf", 1, 2013)
	assert.Equal(t, "remuneracoes-servidores-tjpb-01-2013", serv)
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
func TestSaveFail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	err := save("testFile", ts.URL)
	assert.Error(t, err)
	_, err = os.Stat("testFile.pdf")
	assert.Error(t, err)
}
