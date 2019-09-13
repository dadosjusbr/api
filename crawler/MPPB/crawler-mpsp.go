package main

import (
	"fmt"
	"flag"
	"strings"
	"net/http"
	"io/ioutil"
	"github.com/jackdanger/collectlinks" //Da uma olhada melhor nessa library
)

const linkPrincipal = "http://www.mpsp.mp.br/portal/page/portal/Portal_da_Transparencia/Contracheque/"
var categorias = [7]string{
	"Membros_ativos",
	"Membros_inativos",
	"Servidores_ativos",
	"servidores_inativos",
	"Pensionistas",
	"valores_colaboradores",
	"Verbas-exec-anteriores",
}

func main(){
	mes := flag.Int("mes", 1234, "O mês da planilha")
	ano := flag.Int("ano", 1234, "O ano da planilha")
	
	flag.Parse()

	fmt.Println("Mês: ")
	fmt.Println(*mes)
	fmt.Println("Ano: ")
	fmt.Println(*ano)

	imprimeLinksDePlanilhasOds("http://www.mpsp.mp.br/portal/page/portal/Portal_da_Transparencia/Contracheque/Membros_ativos")
}

func imprimeLinksDePlanilhasOds(url string){
	resp, erro := http.Get(url)

	if erro != nil {
		return
	}

	links := collectlinks.All(resp.Body)  

	for _, link := range(links) {  
		
		if strings.HasSuffix(link, "ods") {
			fmt.Println(link)
		}

	}  
}

func imprimeUmaPagina(url string){
	resp, erro := http.Get(url)

	if erro != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}


