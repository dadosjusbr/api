package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

const (
	urlPrincipal = "http://www.mpm.mp.br/sistemas/consultaFolha/php/RelatorioRemuneracaoMensal.php?"
	extensao     = ".xlsx"
)

var categorias = map[int]string{
	1: "Membros_Ativos",
	2: "Membros_Inativos",
	3: "Servidores_Ativos",
	4: "Servidores_Inativos",
	5: "Pensionistas",
	6: "Colaboradores",
}

func main() {
	mes := flag.Int("mes", -1, "O mês da planilha")
	ano := flag.Int("ano", -1, "O ano da planilha")
	flag.Parse()
	if erro := validaFlags(*mes, *ano); erro != nil {
		fmt.Println(erro)
		return
	}
	urls := geraUrlsDasPlanilhas(*mes, *ano)
	for i, url := range urls {
		arquivo, erro := baixaArquivo(url)
		if erro != nil {
			return
		}
		salvaArquivo(arquivo, *ano, *mes, categorias[i+1])
	}
}

func salvaArquivo(c []byte, ano int, mes int, categoria string) error {
	fileName := fmt.Sprintf("%s-%d-%d%s", categoria, mes, ano, extensao)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Erro ao criar aquivo(%s): %q", fileName, err)
	}
	defer file.Close()
	if _, err = file.Write(c); err != nil {
		return fmt.Errorf("Erro ao escrever arquivo(%s): %q", fileName, err)
	}
	return nil
}

func baixaArquivo(url string) ([]byte, error) {
	body, erro := http.Get(url)
	if erro != nil {
		return nil, fmt.Errorf("erro ao ler URL(%s): %q", url, erro)
	}
	defer body.Body.Close()

	// Tranforma o arquivo em um slice de bytes
	arquivo, erro := ioutil.ReadAll(body.Body)
	if erro != nil {
		return nil, fmt.Errorf("Erro ao fazer parser do arquivo(%s): %q ", url, erro)
	}
	return arquivo, erro
}

// Gera a URL de todos os arquivos de um determinado mês e ano
func geraUrlsDasPlanilhas(mes int, ano int) []string {
	var urls []string
	keys := make([]int, 0)
	for k, _ := range categorias {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, key := range keys {
		url := fmt.Sprintf("%sgrupo=%d&mes=%d&ano=%d", urlPrincipal, key, mes, ano)
		urls = append(urls, url)
	}
	return urls
}

func validaFlags(mes int, ano int) error {
	switch {
	case mes == -1 || ano == -1:
		return errors.New("Mês e ano são flags necessarias")
	case ano > 2019 || ano < 2012:
		return errors.New("Ano inválido")
	case mes > 12 || mes < 1:
		return errors.New("Mês inválido")
	}
	return nil
}
