package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gocolly/colly"
)

// Structure tigerfact type
type TigerFact struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

// Main function
func main() {
	// slice holds facts
	allFacts := make([]TigerFact, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains("factretriever", "www.factretriever.com"),
	)
	// Function allows you to register a callback for when the collector
	// reaches a portion of a page that matches a specific HTML tag specifier
	collector.OnHTML(".factsList li", func(element *colly.HTMLElement) {
		factId, err := strconv.Atoi(element.Attr("id"))
		// returns error
		if err != nil {
			log.Println("Sorry, Could Not Get ID")
		}

		factDetails := element.Text

		fact := TigerFact{
			ID:          factId,
			Description: factDetails,
		}

		allFacts = append(allFacts, fact)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit("https://www.factretriever.com/tiger-facts")

	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", " ")
	// enc.Encode(allFacts)

	writeJSON(allFacts)

}

// returns json encoded and err data
func writeJSON(data []TigerFact) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create JSon file")
		return
	}

	_ = ioutil.WriteFile("tiger-facts.json", file, 0644)
}
