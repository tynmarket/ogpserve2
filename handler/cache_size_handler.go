package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tynmarket/ogpserve2/spider"
)

// CacheSizeHandler return current cache size
type CacheSizeHandler struct {
	Spider *spider.Spider
}

func (h *CacheSizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{}
	//resp["size"] = h.Spider.CurrentCacheSize()

	bytes, err := json.Marshal(resp)
	if err != nil {
		// TODO: Do something
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
