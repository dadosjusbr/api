package crawler

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	basePath        = "http://www.cnj.jus.br"
	remuneracaoPath = "http://www.cnj.jus.br/transparencia/remuneracao-dos-magistrados/remuneracao-"
)

var months = map[int]string{
	1:  "janeiro",
	2:  "fevereiro",
	3:  "marco",
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

// Download download all the spreadsheets related to a month/year and zips it. If successful it returns
// zip file full path.
func Download(month, year int) ([]string, error) {
	if month < 0 || month > 12 {
		return nil, fmt.Errorf("Invalid month: %d", month)
	}
	if year < 2017 {
		return nil, fmt.Errorf("Invalid year: %d", year)
	}
	if year == 2017 && month < 11 {
		return nil, fmt.Errorf("So far the CNJ have not opened data before this nov-2017. You requested: %d/%d", month, year)
	}

	resp, err := http.Get(fmt.Sprintf("%s%s-%d", remuneracaoPath, months[month], year))
	if err != nil {
		return nil, err
	}
	// If the page is not yet there, there is nothing we could do.
	if resp.StatusCode == 404 {
		return []string{}, nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	dir, err := ioutil.TempDir("", "dadosjusbr")
	if err != nil {
		return nil, fmt.Errorf("Error creating temporary directory to store spreadsheets of %d/%d: %q", month, year, err)
	}

	var files []string
	sel := doc.Find("td")
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		if strings.HasSuffix(link, "xls") || strings.HasSuffix(link, "xlsx") {
			dLink := fmt.Sprintf("%s%s", basePath, link)
			resp, err := http.Get(dLink)
			if err != nil {
				return nil, fmt.Errorf("Error making get request (%s): %q", dLink, err)
			}
			// Reading spreadsheet contents.
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				return nil, fmt.Errorf("Error reading response body:%q", err)
			}
			resp.Body.Close()
			outPath := filepath.Join(dir, filepath.Base(dLink))
			if err := writeFile(outPath, contents); err != nil {
				return nil, fmt.Errorf("Error writing file contents: %q", err)
			}
			files = append(files, outPath)
			fmt.Printf("%s downloaded to %s\n", dLink, outPath)
		}
	}
	return files, nil
}

func writeFile(name string, contents []byte) error {
	f, err := os.Create(name)
	defer f.Close()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	if _, err := w.Write(contents); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
