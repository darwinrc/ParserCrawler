package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"server/internal/infra"
	"server/internal/model"
	"server/internal/repo"
)

const UrlNotFound = "url not found in cache"

type CrawlerService interface {
	Crawl(ctx context.Context, url string) (*model.Response, error)
	ConsumeFromRequestQueue(ctx context.Context, broadcast chan []byte)
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

func (s *crawlService) Crawl(ctx context.Context, url string) (*model.Response, error) {
	log.Printf("Crawling URL: %s...\n", url)

	reqId := uuid.New().String()

	data, err := s.CrawlerRepo.GetUrl(ctx, url)
	if err != nil {
		if err.Error() == repo.KeyNotFound {
			go s.publishToRequestQueue(url, reqId)

			return &model.Response{
				Request: model.Request{
					ReqId: reqId,
					Url:   url,
				},
				Response: "",
			}, errors.New(UrlNotFound)
		}

		return nil, errors.New(fmt.Sprintf("error getting url: %s", err))
	}

	fmt.Printf("Data returned from cache: %s...\n", data)

	return &model.Response{
		Request: model.Request{
			ReqId: reqId,
			Url:   url,
		},
		Response: data,
	}, nil
}

func (s *crawlService) publishToRequestQueue(url string, reqId string) {
	log.Println("Publishing url to request queue: ", url)
	req := model.Request{
		Url:   url,
		ReqId: reqId,
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Printf("error marshaling request: %s", err)
		return
	}

	err = s.AMQPClient.SetupAMQExchange()
	if err != nil {
		log.Printf("error setting up the amq connection and exchange: %s", err)
	}

	if err := s.AMQPClient.PublishAMQMessage(body); err != nil {
		log.Printf("error publishing to the exchange: %s", err)
	}

	log.Printf("URL published: %s\n", body)
}

func (s *crawlService) ConsumeFromRequestQueue(ctx context.Context, broadcast chan []byte) {
	err := s.AMQPClient.SetupAMQExchange()
	if err != nil {
		log.Printf("error setting up the amq connection and exchange: %s", err)
	}

	messages, err := s.AMQPClient.ConsumeAMQMessages()
	if err != nil {
		log.Printf("error consuming messages: %s", err)
	}

	res := &model.Response{}

	for msg := range messages {
		if err := json.Unmarshal(msg.Body, res); err != nil {
			log.Printf("error unmarshaling response: %s", err)
			return
		}
		s.CrawlerRepo.StoreUrl(ctx, res.Url, res.Response)

		broadcast <- msg.Body
	}
}
