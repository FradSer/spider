package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// Publication stores information about a publication
type Publication struct {
	Title string
	// Authors  string
	// Abstract string
	URL string
}

func main() {
	fName := "chi_2019.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: sigchi.org, st.sigchi.org
		colly.AllowedDomains("sigchi.org", "st.sigchi.org"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./chi_2019_cache"),
	)

	// courses := make([]Course, 0, 200)
	publications := make([]Publication, 0, 200)

	// On every <a> element which has "href" attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		title := e.Text
		link := e.Attr("href")

		publication := Publication{
			Title: title,
			// Authors:  e.ChildText("li.banner-instructor-info > a > div > div > span"),
			URL: link,
			// Abstract: e.ChildText("div.content"),
		}
		publications = append(publications, publication)

		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.Visit("http://st.sigchi.org/publications/toc/chi-2019-ea.html")

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(publications)
}
