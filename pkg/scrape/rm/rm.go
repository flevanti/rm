package rm

import (
	"fmt"
	"strings"
	"time"
)
import "github.com/gocolly/colly"

var baseUrl = "https://www.rightmove.co.uk/property-for-sale/%s.html?radius=%s&propertyTypes=&mustHave=&dontShow=&furnishTypes=&keywords=&includeSSTC=true&sortType=6&index=%d"
var radius = "1.0"
var pageIndexIncrement = 24
var pageIndexInitial = 0
var coll *colly.Collector
var collInit bool

type PropertyT struct {
	PropertyId string
	Price      int
	Currency   string
	status     string
	status2    string
}

type ScrapeResult struct {
	Results    int
	Properties []PropertyT
	TimeSpent  time.Duration
	Error      error
}

func Postcode(p string) ScrapeResult {
	processStart := time.Now()

	url := buildUrl(p)
	ctx := colly.NewContext()
	ctx.Put("postcode", p)
	ctx.Put("currentPage", "1")
	ctx.Put("lastPage", "0")
	ctx.Put("hasResults", "0")
	getColl().Request("GET", url, nil, ctx, nil)
	return ScrapeResult{TimeSpent: time.Since(processStart)}

}

func buildUrl(p string) string {
	url := fmt.Sprintf(baseUrl, p, radius, pageIndexInitial)
	return url
}

func getColl() *colly.Collector {
	if !collInit {
		initCollector()
	}
	return coll
}

func initCollector() {
	coll = colly.NewCollector()
	coll.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	coll.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	coll.OnHTML("body", checkIfResults)
	coll.OnHTML(".l-searchResult.is-list", processResultsPage)
	// Set error handler
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", string(r.Body), "\nError:", err.Error())
	})
}

func checkIfResults(h *colly.HTMLElement) {
	fmt.Println("checking...")
	if strings.Contains(h.Text, "0 results") {
		h.Request.Ctx.Put("hasResults", "0")
	}
}

func processResultsPage(h *colly.HTMLElement) {

	fmt.Println(h.Attr("id"))
	price := h.DOM.Find("div .propertyCard-priceValue")
	fmt.Println(price.Text())
	added := h.DOM.Find("div .propertyCard-branchSummary")
	fmt.Println(added.Text())
	added = h.DOM.Find(".propertyCard-contactsAddedOrReduced")
	fmt.Println(added.Text())
	fmt.Println("----------------")
}
