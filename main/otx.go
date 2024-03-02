package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	get      = "GET"
	baseurl  = "https://otx.alienvault.com/api/v1/pulses/subscribed/"
	pageSize = 10  // Number of results per page
	maxPages = 100 // Maximum number of pages to retrieve
)

type AVresult struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	AuthorName        string         `json:"author_name"`
	Modified          string         `json:"modified"`
	Created           string         `json:"created"`
	Revision          int            `json:"revision"`
	Tlp               string         `json:"tlp"`
	Public            int            `json:"public"`
	Adversary         string         `json:"adversary"`
	Tags              []string       `json:"tags"`
	TargetedCountries []interface{}  `json:"targeted_countries"`
	MalwareFamilies   []string       `json:"malware_families"`
	AttackIds         []string       `json:"attack_ids"`
	References        []string       `json:"references"`
	Industries        []interface{}  `json:"industries"`
	ExtractSource     []interface{}  `json:"extract_source"`
	MoreIndicators    bool           `json:"more_indicators"`
	Indicators        []OTXIndicator `json:"indicators,omitempty"`
}

type OTXIndicator struct {
	ID          int64       `json:"_id"`
	Indicator   string      `json:"indicator"`
	Type        string      `json:"type"`
	Created     string      `json:"created"`
	Content     string      `json:"content"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description"`
	Expiration  interface{} `json:"expiration"`
	IsActive    int         `json:"is_active"`
	Role        interface{} `json:"role"`
}

type PageInfo struct {
	Count            int         `json:"count"`
	PrefetchPulseIds bool        `json:"prefetch_pulse_ids"`
	T                int         `json:"t"`
	T2               float64     `json:"t2"`
	T3               float64     `json:"t3"`
	Previous         interface{} `json:"previous"`
	Next             string      `json:"next"`
}

type ResponseWithResults struct {
	*http.Response
	Results []AVresult `json:"results,omitempty"`
}

func main() {
	var OTX_API_KEY string
	client := &http.Client{}

	fmt.Println("Please enter your OTX_API_KEY:")
	fmt.Scanln(&OTX_API_KEY)

	domainFile, err := os.Create("domain_indicators.txt")
	if err != nil {
		log.Fatal("Error creating domain file:", err)
	}
	defer domainFile.Close()

	IPv4File, err := os.Create("ipv4_indicators.txt")
	if err != nil {
		log.Fatal("Error creating domain file:", err)
	}
	defer domainFile.Close()

	urlFile, err := os.Create("url_indicators.txt")
	if err != nil {
		log.Fatal("Error creating domain file:", err)
	}
	defer domainFile.Close()

	// Perform the first request to get the initial set of data
	url := fmt.Sprintf("%s?api_key=%s&page=1&page_size=%d", baseurl, OTX_API_KEY, pageSize)
	processPages(client, OTX_API_KEY, url, domainFile, urlFile, IPv4File)
}

func processPages(client *http.Client, apiKey, url string, domainFile *os.File, urlFile *os.File, IPv4File *os.File) {
	for page := 1; page <= maxPages; page++ {
		req, err := http.NewRequest(get, url, nil)
		if err != nil {
			log.Fatal("Error creating request:", err)
		}

		req.Header.Set("X-OTX-API-KEY", apiKey)
		req.Header.Add("Accept", "application/json")

		response, err := client.Do(req)
		if err != nil {
			log.Fatal("Error sending request to server:", err)
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Error reading response body:", err)
		}

		var pageResponse ResponseWithResults
		if err := json.Unmarshal(body, &pageResponse); err != nil {
			log.Fatal("Error unmarshaling JSON:", err)
		}

		for _, result := range pageResponse.Results {
			for _, indicator := range result.Indicators {
				// Check indicator type and write to the corresponding file
				switch indicator.Type {
				case "domain":
					domainFile.WriteString(indicator.Indicator + "\n")
				case "URL":
					urlFile.WriteString(indicator.Indicator + "\n")
				case "IPv4":
					IPv4File.WriteString(indicator.Indicator + "\n")
					// Add more cases for other indicator types as needed
				}
			}
		}

		var pageInfo PageInfo
		if err := json.Unmarshal(body, &pageInfo); err != nil {
			log.Fatal("Error unmarshaling PageInfo:", err)
		}

		if pageInfo.Next == "" {
			// No more pages, break out of the loop
			break
		}

		// Construct the URL for the next page
		url = pageInfo.Next
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
