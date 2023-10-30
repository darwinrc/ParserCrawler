package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"server/internal/infra"
	"server/internal/model"
	"server/internal/repo"
)

type CrawlerService interface {
	Crawl(ctx context.Context, url string) (*model.Response, error)
	ConsumeFromResponseQueue(ctx context.Context, broadcast chan []byte)
}

type crawlService struct {
	CrawlerRepo repo.CrawlerRepo
	AMQPClient  infra.AMQPClient
}

// NewCrawlerService builds a service and injects its dependencies
func NewCrawlerService(crawlerRepo repo.CrawlerRepo, amqpClient infra.AMQPClient) CrawlerService {
	return &crawlService{
		CrawlerRepo: crawlerRepo,
		AMQPClient:  amqpClient,
	}
}

const UrlNotFound = "url not found in cache"

// Crawl gets the crawled url data from the cache.
// If there is a cache miss, it publishes the url to the request queue to be processed by the workers
func (s *crawlService) Crawl(ctx context.Context, url string) (*model.Response, error) {
	reqId := uuid.New().String()

	//data, err := s.CrawlerRepo.GetUrl(ctx, url)
	//if err != nil {
	//	if err.Error() != repo.KeyNotFound {
	//		return nil, errors.New(fmt.Sprintf("error getting url from cache: %s", err))
	//	}
	//
	//	go s.publishToRequestQueue(url, reqId)
	//
	//	return &model.Response{
	//		Request: model.Request{
	//			ReqId: reqId,
	//			Url:   url,
	//		},
	//		Status: "accepted for async processing",
	//	}, errors.New(UrlNotFound)
	//}

	go s.publishToRequestQueue(url, reqId)
	return &model.Response{
		Request: model.Request{
			ReqId: reqId,
			Url:   url,
		},
		Status: "accepted for async processing",
	}, errors.New(UrlNotFound)

	//log.Printf("data returned from cache: %s...\n", data)
	//
	//return &model.Response{
	//	Request: model.Request{
	//		ReqId: reqId,
	//		Url:   url,
	//	},
	//	Response: data,
	//}, nil
}

// publishToRequestQueue publishes the url to the request queue to be processed by the workers
func (s *crawlService) publishToRequestQueue(url string, reqId string) {
	req := model.Request{
		Url:   url,
		ReqId: reqId,
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("error marshaling request: %s", err)
		return
	}

	if err = s.AMQPClient.SetupAMQExchange(); err != nil {
		log.Printf("error setting up the amq connection and exchange: %s", err)
		return
	}

	if err := s.AMQPClient.PublishAMQMessage(body); err != nil {
		log.Printf("error publishing to the exchange: %s", err)
		return
	}

	log.Printf("publishing url to the request queue: %s\n", body)
}

// ConsumeFromResponseQueue consumes messages from the response queue and pushes them to the broadcast channel
func (s *crawlService) ConsumeFromResponseQueue(ctx context.Context, broadcast chan []byte) {
	if err := s.AMQPClient.SetupAMQExchange(); err != nil {
		log.Printf("error setting up the amq connection and exchange: %s", err)
		return
	}

	messages, err := s.AMQPClient.ConsumeAMQMessages()
	if err != nil {
		log.Printf("error consuming messages: %s", err)
		return
	}

	res := &model.Response{}

	for msg := range messages {
		if err := json.Unmarshal(msg.Body, res); err != nil {
			log.Printf("error unmarshaling response: %s", err)
			return
		}

		// Store the URL in the cache before sending it to the broadcast channel
		s.CrawlerRepo.StoreUrl(ctx, res.Url, string(msg.Body))

		broadcast <- msg.Body
	}
}
