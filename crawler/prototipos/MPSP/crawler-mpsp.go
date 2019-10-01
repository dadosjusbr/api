package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const linkPrincipal = "http://www.mpsp.mp.br/portal/page/portal/Portal_da_Transparencia/Contracheque/"

//Last piece of URL where the links are found for each category of table
var categories = []string{
	"Membros_ativos",
	"Membros_inativos",
	"Servidores_ativos",
	"servidores_inativos",
	"valores_colaboradores",
	//Pensionistas separados
	"Pensionistas/Pensionistas_membros",
	"Pensionistas/Pensionistas_servidores",
	//Verbas anteriores separadas
	"Verbas-exec-anteriores/Verbas-exec-anteriores-Membros/Ativos_membros",
	"Verbas-exec-anteriores/Verbas-exec-anteriores-Membros/Inativos_membros",
	"Verbas-exec-anteriores/Verbas-exec-anteriores-Servidores",
	//"Indenizacoes", Não existe
}

var months = map[int]string{
	1:  "janeiro",
	2:  "fevereiro",
	3:  "março",
	4:  "abril",
	5:  "maio",
	6:  "junho",
	7:  "julho",
	8:  "agosto",
	9:  "setembro",
	10: "outubro",
	11: "novembro",
	12: "dezembro",
}

var fileNames = map[string]string{
	"Membros_ativos":        "membros_ativos",
	"Membros_inativos":      "membros_inativos",
	"Servidores_ativos":     "servidores_ativos",
	"servidores_inativos":   "servidores_inativos",
	"valores_colaboradores": "colaboradores",
	//Pensionistas separados
	"Pensionistas/Pensionistas_membros":    "pensionistas_membros",
	"Pensionistas/Pensionistas_servidores": "pensionistas_servidores",
	//Verbas anteriores separadas
	"Verbas-exec-anteriores/Verbas-exec-anteriores-Membros/Ativos_membros":   "exercicios_anteriores_membros_ativos",
	"Verbas-exec-anteriores/Verbas-exec-anteriores-Membros/Inativos_membros": "exercicios_anteriores_membros_inativos",
	"Verbas-exec-anteriores/Verbas-exec-anteriores-Servidores":               "exercicios_anteriores_servidores",
}

// Search for node descendants with the year text inside.
func hasYear(node *goquery.Selection, year string) bool {
	for node.Size() != 0 {
		if strings.Contains(node.Text(), year) {
			return true
		}
		node = node.Children()
	}
	return false
}

// Search for a node in sel containing a given month in it's innerText and return it's href
func contentLink(sel *goquery.Selection, month int) (string, error) {
	monthStr := months[month]
	for i := range sel.Nodes {
		item := sel.Eq(i)
		if strings.Contains(item.Text(), monthStr) {
			link, ok := item.Attr("href")
			if ok {
				return link, nil
			}
			break
		}
	}
	return "", fmt.Errorf("Couldn't find link for month %s", monthStr)
}

func fetchContent(category string, month int, year int) ([]byte, error) {
	link := fmt.Sprintf("%s%s", linkPrincipal, category)
	resCategory, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to make the Get request to (%s)", link)
	}
	defer resCategory.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resCategory.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while casting reader to goquery.document (%s)", link)
	}

	//Select tables of interest. According to html they are inside a td.
	sel := doc.Find("td>table")
	var monthNodes *goquery.Selection
	for i := range sel.Nodes {
		node := sel.Eq(i)
		//Try to find table with year inside, it will be right before the table with the actual links.
		if hasYear(node, fmt.Sprintf("%d", year)) {
			//Nodes with info about the month links.
			monthNodes = sel.Eq(i + 1).Find("a")
			break
		}
	}

	//Verify if there are any nodes in the search
	if monthNodes == nil {
		return nil, fmt.Errorf("Error, couldn't find year %d", year)
	}

	//Find content url for a given month
	ssLink, err := contentLink(monthNodes, month)
	if err != nil {
		return nil, err
	}

	//Fetch content
	resFile, err := http.Get(ssLink)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to make the Get request to (%s)", ssLink)
	}
	defer resFile.Body.Close()

	content, err := ioutil.ReadAll(resFile.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to read content from (%s)", ssLink)
	}
	return content, nil
}

func saveFile(c []byte, year int, month int, category string) error {
	//Create a new file in the cwd
	fileName := fmt.Sprintf("%s-%02d-%d", fileNames[category], month, year)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Error creating file(%s): %q", fileName, err)
	}
	defer file.Close()

	//Write to file
	if _, err = file.Write(c); err != nil {
		return fmt.Errorf("Error writing to file (%s): %q", fileName, err)
	}
	return nil
}

func main() {
	month := flag.Int("mes", -1, "O mês da planilha")
	year := flag.Int("ano", -1, "O ano da planilha")
	flag.Parse()

	if *month == -1 || *year == -1 {
		fmt.Printf("Need flags \"--mes=int --ano=int\"\n")
		return
	}

	for _, category := range categories {
		c, err := fetchContent(category, *month, *year)
		if err != nil {
			fmt.Printf("Error while trying to fetch content (%s %02d-%d): %q\n", fileNames[category], *month, *year, err)
			continue
		}
		err = saveFile(c, *year, *month, category)
		if err = saveFile(c, *year, *month, category); err != nil {
			fmt.Printf("Error saving content to file (%s %02d-%d): %q\n", fileNames[category], *month, *year, err)
		}
	}
}
