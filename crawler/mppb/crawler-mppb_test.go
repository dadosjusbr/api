package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinks(t *testing.T) {
	baseURL := "http://teste/"
	l := links(baseURL, 1, 12)
	assert.Equal(t, 8, len(l), "")
	for typ, u := range l {
		assert.Truef(t, strings.HasPrefix(u, baseURL), "URL base inválida para tipo:%s", typ)
		assert.Truef(t, strings.Contains(u, fmt.Sprintf("mes=%d", 1)), "Mês inválido para tipo:%s url:%s", typ, u)
		assert.Truef(t, strings.Contains(u, fmt.Sprintf("exercicio=%d", 12)), "Ano inválido para tipo:%s url:%s", typ, u)
	}
	for typ, id := range tipos {
		assert.Truef(t, strings.Contains(l[typ], fmt.Sprintf("tipo=%d", id)), "ID do tipo inválido para tipo:%s", typ)
	}
}

func TestDownload(t *testing.T) {
	// O foco desse teste é checar se o resultado da chamada ao servidor é realmente
	// salva no buffer.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	}))
	defer ts.Close()

	// Baixando os conteúdo da resposta da chamada ao servidor de teste para o buffer.
	var buf bytes.Buffer
	assert.NoError(t, download(ts.URL, &buf))
	assert.Equal(t, "Hello", buf.String()) //Checando se o que tem no buffer é igual ao retorno da chamada ao servidor.
}
