package main

import{
	"fmt"
	"net/http"
	"strings"
	"time"
	"math/rand"
	"net/url"
	"github.com/PuerkitoBio/goquery"
}

var googleDomains = map[string]string{

}

type SearchResult struct{
	ResultRank int
	ResultURL string
	ResultTitle string
	ResultDesc string
}

var userAgents = []string{

}

func randomUserAgent() string{
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func buildGoogleUrls()(searchTerm){
	toScrape := []string{}
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
}

func GoogleScrape(searchTerm)([]SearchResult, err){
	result := []SearchResult{}
	resultCounter := 0
	googlePages, err := buildGoogleUrls(searchTerm)
}

func main(){
	res, err := GoogleScrape("sarah rubi")
	if err == nil{
		for _, res := range res{
			fmt.Println(res)
		}
	}
}