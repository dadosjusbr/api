package crawler

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const basePath = "http://www.cnj.jus.br"

//Result represents an downloaded spreadsheet with the file name and its bytes.
type Result struct {
	Name string
	Body []byte
}

//Results is an array of Result
type Results []Result

var transport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 60 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 30 * time.Second,
}
var client = &http.Client{
	Timeout:   60 * time.Second,
	Transport: transport,
}

func init() {
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
}

// Crawl download all the spreadsheets related to the page of the given url. If successful it returns
// an array of results with the name and the bytes of each spreadsheet.
func Crawl(url string) (Results, error) {
	resp, err := client.Get(url)
	if err != nil {
		return Results{}, err
	}
	// If the page is not yet there, there is nothing we could do.
	if resp.StatusCode == 404 {
		return Results{}, nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Results{}, err
	}
	results := Results{}
	sel := doc.Find("a")
	for i := range sel.Nodes {
		item := sel.Eq(i)
		link, _ := item.Attr("href")
		if strings.HasSuffix(link, "xls") || strings.HasSuffix(link, "xlsx") {
			var dLink string
			if strings.HasPrefix(link, "file://") {
				dLink = link
			} else {
				dLink = fmt.Sprintf("%s%s", basePath, link)
			}
			resp, err := client.Get(dLink)
			if err != nil {
				return Results{}, fmt.Errorf("Error making get request (%s): %q", dLink, err)
			}
			defer resp.Body.Close()
			// Reading spreadsheet contents.
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return Results{}, fmt.Errorf("Error reading response body:%q", err)
			}
			fn := filepath.Base(link)
			result := Result{fn, contents}
			results = append(results, result)
			fmt.Printf("%s downloaded\n", fn)
		}
	}
	return results, nil
}
