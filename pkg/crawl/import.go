package crawl

import "fmt"
import "github.com/gocolly/colly"

var baseUrl = "https://www.rightmove.co.uk/property-for-sale/%s.html?radius=%s&propertyTypes=&mustHave=&dontShow=&furnishTypes=&keywords=&includeSSTC=true&sortType=6&index=%d"
var radius = "1.0"
var pageIndexIncrement = 24
var pageIndexInitial = 0
var coll *colly.Collector
var collInit bool

func Postcode(p string) {
	url := buildUrl(p)
	getColl().Visit(url)
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
	coll.OnHTML(".l-searchResult.is-list", processResultsPage)
}

func processResultsPage(h *colly.HTMLElement) {
	fmt.Println(h.Index)
	fmt.Println(h.Attr("id"))
	price := h.DOM.Find("div .propertyCard-priceValue")
	fmt.Println(price.Text())
	added := h.DOM.Find("div .propertyCard-branchSummary")
	fmt.Println(added.Text())
	added = h.DOM.Find(".propertyCard-contactsAddedOrReduced")
	fmt.Println(added.Text())

	fmt.Println("----------------")
}
