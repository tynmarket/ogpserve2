package main

import (
	"log"
	"net/http"

	"github.com/tynmarket/ogpserve2/handler"
	"github.com/tynmarket/ogpserve2/spider"
)

func main() {
	spider := &spider.Spider{}

	go spider.Loop()

	// Return Twitter Card
	http.Handle("/twitter", &handler.TwitterCardHandler{Spider: spider})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
