package main

import (
	"flag"
	"fmt"
	"strconv"
)

const linkPrincipal = "http://www.mpsp.mp.br/portal/page/portal/Portal_da_Transparencia/Contracheque/"

var categorias = []string{
	"/Membros_ativos/Tabela%20I%20membros%20ativos%20ref",
	"Membros_inativos/Tabela%20I%20membros%20inativos%20ref",
	"/Servidores_ativos/Tabela%20I%20%20servidores%20ativos%20ref",
	"/servidores_inativos/Tabela%20I%20servidores%20inativos%20ref",
	"Pensionistas/Pensionistas_membros/Benefici%C3%A1rios%20Membros%20ref%",
	"/Pensionistas/Pensionistas_servidores/Benefici%C3%A1rios%20Servidores%20ref%",
	//"/valores_colaboradores/Tabela%20Portal%20colaboradores%20jan%202019.pdf",  TODO: Preciso olhar essa categoria!!
	"/Verbas-exec-anteriores/Verbas-exec-anteriores-Servidores/Tabela%20II%20servidores%20ativos%20ref",
	"/Verbas-exec-anteriores/Verbas-exec-anteriores-Servidores/Tabela%20II%20servidores%20ativos%20ref",
	"/Verbas-exec-anteriores/Verbas-exec-anteriores-Membros/Inativos_membros/Tabela%20II%20membros%20inativos%20ref",
}

func main() {
	mes := flag.String("mes", "", "O mês da planilha")
	ano := flag.Int("ano", -1, "O ano da planilha")
	flag.Parse()

	msgErro, _ := validaFlags(*mes, *ano)
	if msgErro != "" {
		fmt.Println(msgErro)
		return
	}

	fmt.Printf("Mês: %s \n", *mes)
	fmt.Printf("Ano: %d \n", *ano)
	imprimeLinksDeUmAnoMesEspecifico(*mes, *ano)
}

func validaFlags(mes string, ano int) (string, error) {
	mesStr, erro := strconv.Atoi(mes)
	if erro != nil {
		return "", erro
	}
	if mes == "" || ano == -1 {
		return "As flags --mes e --ano são necessarias", erro
	} else if ano > 2019 || ano < 0 {
		return "Não é possivel processar esse ano", erro
	} else if mesStr > 12 || mesStr < 1 {
		return "Mês invalido", erro

	} else {
		return "", erro
	}
}

func imprimeLinksDeUmAnoMesEspecifico(mes string, ano int) {
	for _, categoria := range categorias {
		switch ano {
		case 2018:
			imprime2018(mes, categoria)
		case 2019:
			imprime2019(mes, categoria)
		default:
			fmt.Println("Não implementado para esse ano ainda")
		}
	}
}

func imprime2018(mes string, categoria string) {
	var url string
	if mes == "12" {
		url = linkPrincipal + "Membros_ativos/Tabela%20I%20membros%20ativos%20ref122018.ods" //O mês 12 é o unico mês que foge do padrão no ano de 2018.
	} else {
		url = linkPrincipal + categoria + mes + ".ods"
	}
	fmt.Printf("Categoria: %s \n", categoria)
	fmt.Printf("link: %s \n", url)
}

func imprime2019(mes string, categoria string) {
	url := linkPrincipal + categoria + mes + "19.ods"
	fmt.Printf("Categoria: %s \n", categoria)
	fmt.Printf("link: %s \n", url)
}
