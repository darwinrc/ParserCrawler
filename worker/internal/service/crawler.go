package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"worker/internal/infra"
)

type CrawlerService interface {
	CrawlURL()
}

type crawlService struct {
	AMQPClient infra.AMQPClient
}

// NewCrawlerService builds a service and injects its dependencies
func NewCrawlerService(amqpClient infra.AMQPClient) CrawlerService {
	return &crawlService{
		AMQPClient: amqpClient,
	}
}

type Request struct {
	ReqId string `json:"reqId,omitempty"`
	Url   string `json:"url,omitempty"`
}

type Response struct {
	Request
	Response string `json:"response"`
}

func (s *crawlService) CrawlURL() {
	fmt.Println("Crawling URL...")

	err := s.AMQPClient.SetupAMQExchange()
	if err != nil {
		log.Printf("error setting up the amq connection and exchange: %s", err)
	}

	messages, err := s.AMQPClient.ConsumeAMQMessages()
	if err != nil {
		log.Printf("error consuming messages: %s", err)
	}

	for msg := range messages {
		req := &Request{}
		if err := json.Unmarshal(msg.Body, req); err != nil {
			log.Printf("error unmarshaling payload: %s", err)
			return
		}

		data := s.crawl(req.Url)

		res := &Response{
			Request: Request{
				ReqId: req.ReqId,
				Url:   req.Url,
			},
			//ReqId:    req.ReqId,
			//Url:      req.Url,
			Response: data,
		}

		body, err := json.Marshal(res)
		if err != nil {
			log.Printf("error marshaling payload: %s", err)
			return
		}

		if err := s.AMQPClient.PublishAMQMessage(body); err != nil {
			log.Printf("error publishing to the exchange: %s", err)
		}

		log.Printf("Published to the exchange: %#v", string(body))
	}

}

func (s *crawlService) crawl(url string) string {
	log.Printf("About to crawl url: %#v\n", url)

	time.Sleep(5 * time.Second)
	return `crawled: ` + url
}

//var (
//	visitedURLs = make(map[string]bool)
//	visitedMu   sync.Mutex
//)

//func main() {
//
//	//seedURL := "https://parserdigital.com/" // Replace with your starting URL
//	//domain := "https://parserdigital.com"   // Replace with your target subdomain
//	//robotsFile := "robots.txt"
//	//
//	//delay := getDelay(seedURL, robotsFile)
//
//	//crawl(seedURL, domain, delay)
//
//}

//func getDelay(urlStr, robotsFile string) int {
//	resp, err := http.Get(fmt.Sprintf("%s%s", urlStr, robotsFile))
//	if err != nil {
//		fmt.Println("Error:", err)
//		return 0
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return 0
//	}
//
//	lines := strings.Split(string(body), "\n")
//	for _, line := range lines {
//		parts := strings.SplitN(line, ":", 2)
//		if len(parts) != 2 || !strings.EqualFold(strings.TrimSpace(parts[0]), "Crawl-delay") {
//			continue
//		}
//
//		delay, err := strconv.Atoi(strings.TrimSpace(parts[1]))
//		if err != nil {
//			fmt.Println("Error:", err)
//			return 0
//		}
//
//		return delay
//
//	}
//
//	return 0
//}
//
//func crawl(urlStr, targetSubdomain string, delay int) {
//	// Ensure we only visit URLs within the target subdomain
//	if !strings.Contains(urlStr, targetSubdomain) {
//		return
//	}
//
//	// Mark the URL as visited to avoid duplicates
//	visitedMu.Lock()
//	visitedURLs[urlStr] = true
//	visitedMu.Unlock()
//
//	// Wait for crawl delay from robots.txt
//	time.Sleep(time.Duration(delay) * time.Second)
//
//	// Make an HTTP request to the URL
//	resp, err := http.Get(urlStr)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	// Parse the HTML content using goquery
//	doc, err := goquery.NewDocumentFromReader(resp.Body)
//	if err != nil {
//		fmt.Println("Error parsing document:", err)
//		return
//	}
//
//	// Extract and print links on the current page
//	fmt.Println("Visited URL:", urlStr)
//	doc.Find("a").Each(func(index int, element *goquery.Selection) {
//		href, exists := element.Attr("href")
//		if exists {
//			// Resolve relative URLs to absolute URLs
//			linkURL, err := url.Parse(href)
//			if err == nil {
//				absoluteURL := resp.Request.URL.ResolveReference(linkURL)
//				link := absoluteURL.String()
//
//				// Ensure the link belongs to the same subdomain
//				if strings.Contains(link, targetSubdomain) && link != urlStr {
//					fmt.Println("Found link:", link)
//				}
//			}
//		}
//	})
//
//	// Extract and process links in parallel
//	var wg sync.WaitGroup
//	doc.Find("a").Each(func(index int, element *goquery.Selection) {
//		href, exists := element.Attr("href")
//		if exists {
//			// Resolve relative URLs to absolute URLs
//			linkURL, err := url.Parse(href)
//			if err == nil {
//				absoluteURL := resp.Request.URL.ResolveReference(linkURL)
//				link := absoluteURL.String()
//
//				// Ensure the link belongs to the same subdomain
//				if strings.Contains(link, targetSubdomain) {
//					wg.Add(1)
//					go func(linkURL string) {
//						defer wg.Done()
//						if !visited(linkURL) {
//							crawl(linkURL, targetSubdomain, delay)
//						}
//					}(link)
//				}
//			}
//		}
//	})
//
//	// Wait for all Goroutines to finish
//	wg.Wait()
//}
//
//func visited(urlStr string) bool {
//	visitedMu.Lock()
//	defer visitedMu.Unlock()
//	return visitedURLs[urlStr]
//}
