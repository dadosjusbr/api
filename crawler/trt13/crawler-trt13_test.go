package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestDownloadFail(t *testing.T) {
	// O foco desse teste é checar se o erro é lançado quando o status code é diferente de 200.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	// Baixando os conteúdo da resposta da chamada ao servidor de teste para o buffer.
	var buf bytes.Buffer
	assert.Error(t, download(ts.URL, &buf)) //Função deve lançar um erro para status diferentes de 200...
}
