package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tynmarket/ogpserve2/model"
	"github.com/tynmarket/ogpserve2/spider"
)

// TwitterCardHandler return Twitter Card response
type TwitterCardHandler struct {
	Spider *spider.Spider
}

// ServeHTTP handler response
func (h *TwitterCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ogps := h.Spider.Run(query)

	resp := map[string]model.TwitterCard{}

	for _, ogp := range ogps {
		card := *ogp.MergeIntoTwitter()
		if card.ValuePresent() {
			resp[ogp.RequestURL] = card
		}
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		// TODO: Do something
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
