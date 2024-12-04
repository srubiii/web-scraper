package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var domainsGoogle = map[string]string{
	"com": "https://www.google.com/search?q=",
	"ca":  "https://www.google.ca/search?q=",
	"uk":  "https://www.google.co.uk/search?q=",
}

// got domains from https://yourreputations.com/country-specific-google-domains-list/
type InqueryOutput struct {
	OutputOrder int
	OutputURL   string
	OutputTitle string
	OutputDesc  string
}

var agentUsers = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

//got agents from https://useragents.io/

func randomUserAgent() string {
	x := rand.New(rand.NewSource(time.Now().Unix()))
	numRand := x.Int() % len(agentUsers)
	return agentUsers[numRand]
}

func buildUrls(termSearch, codeCountry, languageCode string, pages, count int) ([]string, error) {
	goScrape := []string{}
	termSearch = strings.Trim(termSearch, " ")
	termSearch = strings.Replace(termSearch, " ", "+", -1)
	if base, found := domainsGoogle[codeCountry]; found {
		for i := 0; i < pages; i++ {
			start := i * count
			urlScrape := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0", base, termSearch, count, languageCode, start)
			goScrape = append(goScrape, urlScrape)
		}
	} else {
		err := fmt.Errorf("country (%s) is not currently supported", codeCountry)
		return nil, err
	}
	return goScrape, nil
}

func ResultParsing(communication *http.Response, order int) ([]InqueryOutput, error) {

	doc, err := goquery.NewDocumentFromReader(communication.Body)
	if err != nil {
		return nil, err
	}
	defer communication.Body.Close()

	outputs := []InqueryOutput{}
	sel := doc.Find("div.g")
	order++

	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3.x")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")

		fmt.Printf("Parsed Result - Link: %s, Title: %s, Desc: %s\n", link, title, desc)

		if link != "" && link != "#" && !strings.HasPrefix(link, "/") {
			output := InqueryOutput{
				order,
				link,
				title,
				desc,
			}
			outputs = append(outputs, output)
			order++
		}
	}

	return outputs, err
}

func ClientScrapeGot(stringProxy interface{}) *http.Client {
	switch v := stringProxy.(type) {
	case string:
		proxyUrl, _ := url.Parse(v)
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	default:
		return &http.Client{}
	}
}

func TheScrape(termSearch, codeCountry, languageCode string, stringProxy interface{}, pages, count, backoff int) ([]InqueryOutput, error) {
	outputs := []InqueryOutput{}
	resultCounter := 0
	googlePages, err := buildUrls(termSearch, codeCountry, languageCode, pages, count)
	if err != nil {
		return nil, err
	}
	for _, page := range googlePages {
		res, err := clientRequestScrape(page, stringProxy)
		if err != nil {
			return nil, err
		}
		data, err := ResultParsing(res, resultCounter)
		if err != nil {
			return nil, err
		}
		resultCounter += len(data)
		for _, output := range data {
			outputs = append(outputs, output)
		}
		time.Sleep(time.Duration(backoff) * time.Second)
	}
	return outputs, nil
}

func clientRequestScrape(searchURL string, stringProxy interface{}) (*http.Response, error) {
	clientBase := ClientScrapeGot(stringProxy)
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := clientBase.Do(req)
	if res.StatusCode != 200 {
		err := fmt.Errorf("scraper received a non-200 status code suggesting a ban")
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return res, nil
}

func main() {
	res, err := TheScrape("sarah rubi", "com", "en", nil, 1, 30, 30)
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	} else {
		for _, output := range res {
			fmt.Println(output)
		}
	}
}
