package main

import (
	"net/http"

	"github.com/tynmarket/ogpserve2/handler"
	"github.com/tynmarket/ogpserve2/spider"
)

func main() {
	spider := &spider.Spider{}

	// Return Twitter Card
	http.Handle("/twitter", &handler.TwitterCardHandler{Spider: spider})
}
