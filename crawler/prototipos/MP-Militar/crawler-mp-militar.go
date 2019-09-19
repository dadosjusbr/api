package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const urlPrincipal = "http://www.mpm.mp.br/sistemas/consultaFolha/php/RelatorioRemuneracaoMensal.php?"
const extensao = ".xlsx"

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
	erro := validaFlags(*mes, *ano)
	if erro != nil {
		fmt.Println(erro)
		return
	}
	urls := geraUrlsDasPlanilhas(*mes, *ano)
	for url := range urls {
		arquivo, erro := baixaArquivo(urls[url])
		if erro != nil {
			return
		}
		salvaArquivo(arquivo, *ano, *mes, categorias[url])
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
		return nil, errors.New("Erro ao ler url") //TODO:colocar nome da url
	}
	defer body.Body.Close()
	arquivo, erro := transformaBodyEmSlices(body)
	if erro != nil {
		return nil, errors.New("Erro fazer parser do arquivo") //TODO:colocar nome da url
	}
	return arquivo, erro
}

func transformaBodyEmSlices(response *http.Response) ([]byte, error) {
	targetBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("Erro ao converter arquivo em slice de bytes")
	}
	return targetBody, nil
}

func geraUrlsDasPlanilhas(mes int, ano int) []string {
	var urls []string
	for j := 1; j <= 6; j++ {
		url := fmt.Sprintf("%sgrupo=%d&mes=%d&ano=%d", urlPrincipal, j, mes, ano)
		urls = append(urls, url)
	}
	return urls
}

func validaFlags(mes int, ano int) error {
	if mes == -1 || ano == -1 {
		return errors.New("Mês e ano são flags necessarias")
	} else if ano > 2019 || ano < 2012 {
		return errors.New(("Ano inválido"))
	} else if mes > 12 || mes < 1 {
		return errors.New("Mês inválido")
	}
	return nil
}
