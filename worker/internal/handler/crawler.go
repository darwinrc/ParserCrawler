package handler

import (
	"encoding/json"
	"log"
	"worker/internal/infra"
	"worker/internal/model"
	"worker/internal/service"
)

type CrawlerHandler interface {
	Process()
}

type crawlHandler struct {
	AMQPClient infra.AMQPClient
	Service    service.CrawlerService
}

// NewCrawlerHandler builds a service and injects its dependencies
func NewCrawlerHandler(amqpClient infra.AMQPClient, service service.CrawlerService) CrawlerHandler {
	return &crawlHandler{
		AMQPClient: amqpClient,
		Service:    service,
	}
}

// Process consumes the url to process from the request queue, calls the service to crawl it and publishes the results
// to the response queue
func (h *crawlHandler) Process() {
	err := h.AMQPClient.SetupAMQExchange()
	if err != nil {
		log.Printf("error setting up the amq connection and exchange: %s", err)
		return
	}

	messages, err := h.AMQPClient.ConsumeAMQMessages()
	if err != nil {
		log.Printf("error consuming messages: %s", err)
		return
	}

	for msg := range messages {
		req := &model.Request{}
		if err := json.Unmarshal(msg.Body, req); err != nil {
			log.Printf("error unmarshaling payload: %s", err)
			return
		}

		log.Printf("getting message from the request queue: %v", req)

		data := h.Service.Crawl(req.Url)

		res := &model.Response{
			Request: model.Request{
				ReqId: req.ReqId,
				Url:   req.Url,
			},
			Response: *data,
		}

		body, err := json.Marshal(res)
		if err != nil {
			log.Printf("error marshaling payload: %s", err)
			return
		}

		if err := h.AMQPClient.PublishAMQMessage(body); err != nil {
			log.Printf("error publishing to the exchange: %s", err)
			return
		}

		log.Printf("published to the exchange: %#v", string(body))
	}
}
