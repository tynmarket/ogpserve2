package main

import (
	"net/http"

	"github.com/tynmarket/ogpserve2/handler"
	"github.com/tynmarket/ogpserve2/spider"
)

var QueueSize = 10000

func main() {
	queue := make(chan string, QueueSize)
	spider := &spider.Spider{queue: queue}

	// Return Twitter Card
	http.Handle("/twitter", &handler.TwitterCardHandler{Spider: spider})
}
