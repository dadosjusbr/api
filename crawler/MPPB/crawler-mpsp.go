package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jackdanger/collectlinks" //Da uma olhada melhor nessa library
)

const linkPrincipal = "http://www.mpsp.mp.br/portal/page/portal/Portal_da_Transparencia/Contracheque/"

var categorias = []string{
	"Membros_ativos",
	"Membros_inativos",
	"Servidores_ativos",
	"servidores_inativos",
	"Pensionistas/Pensionistas_membros",
	"Pensionistas/Pensionistas_servidores",
	"valores_colaboradores",
	"Verbas-exec-anteriores/Verbas-exec-anteriores-Servidores",
	"Verbas-exec-anteriores/Verbas-exec-anteriores-Membros/Ativos_membros",
	"Verbas-exec-anteriores/Verbas-exec-anteriores-Membros/Inativos_membros",
}

func main() {
	mes := flag.String("mes", "1234", "O mês da planilha")
	ano := flag.Int("ano", -1, "O ano da planilha")
	flag.Parse()
	fmt.Printf("Mês: %s \n", *mes)
	fmt.Printf("Ano: %d \n", *ano)
	imprimeLinksDeUmAnoMesEspecifico(*mes, *ano)
}

func imprimeLinksDeUmAnoMesEspecifico(mes string, ano int) {
	switch ano {
	case 2018:
		imprime2018(mes)
	case 2019:
		imprime2019(mes)
	default:
		fmt.Println("not implemented yet")
	}
}

func imprime2018(mes string) {
	var url string
	if mes == "12" {
		url = linkPrincipal + "Membros_ativos/Tabela%20I%20membros%20ativos%20ref122018.ods"
	} else {
		url = linkPrincipal + "Membros_ativos/Tabela%20I%20membros%20ativos%20ref" + mes + ".ods"
	}
	fmt.Println(url)
}

func imprime2019(mes string) {
	url := linkPrincipal + "/Membros_ativos/Tabela%20I%20membros%20ativos%20ref" + mes + "19.ods"
	fmt.Println(url)
}

//Imprime os links de todas as planilhas ODS de todas as categorias
func linksDePlanilhasPorCategorias() {
	for _, categoria := range categorias {
		url := linkPrincipal + categoria
		fmt.Println("Categoria: ", categoria)
		imprimeLinksDePlanilhasOds(url)
	}
}

//Imprime os links de todas as planilhas ODS de uma pagina com a url passada como parametro.
func imprimeLinksDePlanilhasOds(url string) error {
	resp, erro := http.Get(url)
	if erro != nil {
		return erro
	}
	links := collectlinks.All(resp.Body)
	for _, link := range links {
		if strings.HasSuffix(link, "ods") {
			fmt.Println(link)
		}
	}
	return erro
}

//Imprime todo o conteudo de uma pagina com a url passada como parametro.
func imprimeUmaPagina(url string) error {
	resp, erro := http.Get(url)
	if erro != nil {
		return erro
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return erro
}
