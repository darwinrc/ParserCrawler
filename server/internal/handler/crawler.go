package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"server/internal/model"
	"server/internal/service"
)

type CrawlerHandler interface {
	Attach(r *mux.Router)
	HandleCrawl(w http.ResponseWriter, r *http.Request)
}

type crawlerHandler struct {
	Service service.CrawlerService
}

// NewCrawlerHandler builds a handler and injects its dependencies
func NewCrawlerHandler(s service.CrawlerService) CrawlerHandler {
	return &crawlerHandler{
		Service: s,
	}
}

// Attach attaches the crawler endpoints to the router
func (h *crawlerHandler) Attach(r *mux.Router) {
	r.HandleFunc("/crawl", h.HandleCrawl).Methods("GET", "OPTIONS")
}

// HandleCrawl exposes the API to crawl a website
func (h *crawlerHandler) HandleCrawl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	url := r.URL.Query().Get("url")

	res, err := h.Service.Crawl(r.Context(), url)
	if err != nil {
		if err.Error() == service.UrlNotFound {
			res = &model.Response{
				Status: "accepted",
			}
			w.WriteHeader(http.StatusAccepted)
		} else {
			res = &model.Response{
				Status: "error",
			}
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}

	body, err := json.Marshal(res)
	if err != nil {
		log.Printf("error marshaling payload: %s", err)
		return
	}

	w.Write(body)
}
