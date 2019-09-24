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
var categories = [...]string{
	"Membros_ativos",
	"Membros_inativos",
	"Servidores_ativos",
	"servidores_inativos",
	"Pensionistas",
	"valores_colaboradores",
	"Verbas_exec_anteriores",
	"Indenizacoes",
}

var meses = map[int]string{
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

// Search for node descendants with the year text inside.
func findYear(node *goquery.Selection, year string) bool {
	for node.Size() != 0 {
		if strings.Contains(node.Text(), year) {
			return true
		}
		node = node.Children()
	}
	return false
}

func fetchContent(category string, month int, year int) ([]byte, error) {
	link := fmt.Sprintf("%s%s", linkPrincipal, category)
	res, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to make the Get request to (%s)", link)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while casting reader to goquery.document (%s)", link)
	}
	res.Body.Close()

	//Select tables of interest. According to html they are inside a td.
	sel := doc.Find("td>table")
	var monthNodes *goquery.Selection
	for i := range sel.Nodes {
		node := sel.Eq(i)
		//Try to find table with year inside, it will be right before the table with the actual links.
		if findYear(node, fmt.Sprintf("%d", year)) {
			//Nodes with info about the month links.
			monthNodes = sel.Eq(i + 1).Find("a")
			break
		}
	}
	//Verify if given month exists
	if monthNodes == nil {
		return nil, fmt.Errorf("Error, couldn't find year %d (%s)", year, link)
	} else if monthNodes.Size() < month-1 {
		return nil, fmt.Errorf("Error, couldn't find month %d (%s)", month, link)
	}
	//Find content url
	ssLink, ok := monthNodes.Eq(month - 1).Attr("href")
	if !ok {
		return nil, fmt.Errorf("Error, couldn't find href for month %d (%s)", month, link)
	}

	//Fetch spreedsheet
	res, err = http.Get(ssLink)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to make the Get request to (%s)", link)
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to read content from (%s)", ssLink)
	}
	return content, nil
}

func saveFile(c []byte, year int, month int, category string) error {
	//Create a new file in the cwd
	fileName := fmt.Sprintf("%s-%02d-%d", category, month, year)
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

	for _, category := range categories {
		c, err := fetchContent("Membros_ativos", *month, *year)
		if err != nil {
			fmt.Printf("Error while trying to fetch content %s %d-%d: %q\n", category, *month, *year, err)
			continue
		}
		err = saveFile(c, *year, *month, category)
		if err = saveFile(c, *year, *month, category); err != nil {
			fmt.Printf("Error saving spreedsheet to file (%s %d-%d): %q\n", category, *month, *year, err)
		}
	}

}
