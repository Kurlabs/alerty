package robotrun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/gocolly/colly"
)

// Attributes is the message data structure
type Attributes struct {
	BrainURL string `json:"brainurl"`
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data       []byte     `json:"data"`
	Attributes Attributes `json:"attributes"`
}

// RResult is an temporal Field to build a Response
type RResult struct {
	URL          string `json:"url"`
	ResponseCode int    `json:"response_code"`
}

// Response is a object related with checks
type Response struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	RobotResult string `json:"robot_result"`
}

// Websites is the basic document to scan
type Websites struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// rresultsJSON transform responses to a json array
func rresultsJSON(rresults []RResult) []byte {
	output, err := json.Marshal(rresults)
	if err != nil {
		return []byte("[]")
	}
	return output
}

// responsesJSON transform responses to a json array
func responsesJSON(responses []Response) []byte {
	output, err := json.Marshal(responses)
	if err != nil {
		return []byte("[]")
	}
	return output
}

func RobotRun(ctx context.Context, m PubSubMessage) error {
	log.Printf("Robot Websites: %s!", m.Data)
	var webs []Websites
	err := json.Unmarshal(m.Data, &webs)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BrainURL: %s", m.Attributes.BrainURL)
	token := os.Getenv("TOKEN")
	brainURL := m.Attributes.BrainURL

	// var webs = []Websites{
	// 	{URL: "https://alerty.online"},
	// 	// {URL: "https://platzi.com/"},
	// 	{URL: "https://about.gitlab.com/"},
	// }
	var responses []Response

	var wg sync.WaitGroup
	wg.Add(len(webs))

	for _, web := range webs {
		go func(URL string, id string) {
			defer wg.Done()
			var rresults []RResult
			var domain string

			re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?([^:\/\n]+)`)

			match := re.FindStringSubmatch(URL)
			if len(match) == 2 {
				domain = match[1]
			}
			re2 := regexp.MustCompile(`[@\?]`)
			c := colly.NewCollector(
				// Visit only domains: Domain
				colly.AllowedDomains(domain),
				colly.DisallowedURLFilters(re2),
				//colly.MaxDepth(5),
				colly.Async(true),

				// Cache responses to prevent multiple download of pages
				// even if the collector is restarted
				// colly.CacheDir("./cache"),
			)
			c.ParseHTTPErrorResponse = true
			c.UserAgent = "Alerty"
			c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

			// On every a element which has href attribute call callback
			c.OnHTML("a[href]", func(e *colly.HTMLElement) {
				link := e.Attr("href")
				// Visit link found on page
				// Only those links are visited which are in AllowedDomains
				c.Visit(e.Request.AbsoluteURL(link))
			})

			// Before making a request print "Visiting ..."
			c.OnRequest(func(r *colly.Request) {
				// fmt.Println("Visiting", r.URL.String())
			})

			c.OnResponse(func(r *colly.Response) {
				fmt.Println(r.Request.URL.String(), r.StatusCode)
				rresults = append(rresults, RResult{URL: r.Request.URL.String(), ResponseCode: r.StatusCode})
			})

			// Start scraping on URL
			c.Visit(URL)
			c.Wait()

			var results = rresultsJSON(rresults)
			responses = append(responses, Response{ID: id, URL: URL, RobotResult: string(results)})
		}(web.URL, web.ID)
	}
	wg.Wait()

	var result = responsesJSON(responses)
	fmt.Println(string(result))
	// Send to API
	if brainURL != "" {
		// Create request object.
		req, err := http.NewRequest("POST", brainURL, bytes.NewBuffer(result))

		// Set the header in the request.
		req.Header.Set("TOKEN", token)

		// Execute the request.
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Fatal("Error connection to Brain API", err)
		}
		defer resp.Body.Close()
	}
	return nil
}
