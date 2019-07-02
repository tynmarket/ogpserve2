package spider

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/tynmarket/ogpserve2/model"
	"golang.org/x/time/rate"
)

// Spider is
type Spider struct {
}

const (
	// クロールの並列数の最大値
	crawlerCount = 10
	logFile      = "/var/log/ogpserve2/ogpserve2.log"
)

var queueSize = 10000
var queue = make(chan string, queueSize)
var cacheSize = 1000
var cache, _ = lru.New(cacheSize)

//var logger = GetLogger()

// Run is
func (s *Spider) Run(query url.Values) []*model.Ogp {
	urls := query["urls[]"]
	//tag := query.Get("tag")
	skipCrawl := query.Get("skip_crawl") == "true"

	if urls == nil || len(urls) == 0 {
		return nil
	}

	//logCount("request with urls", len(urls), tag)

	for _, url := range urls {
		logURL("request for", url)
	}

	// Serve from cache if present
	ogps, urls := s.serve(urls)

	// Read cache only
	if skipCrawl {
		return ogps
	}

	// Crawl for not cached URLs
	for _, url := range urls {
		queue <- url
	}

	return ogps
}

// Loop is
func (s *Spider) Loop() {
	ctx := context.Background()
	n := rate.Every(time.Second / crawlerCount)
	l := rate.NewLimiter(n, crawlerCount)

	domains := make(map[string]int64)
	mutex := new(sync.Mutex)
	parser := &Parser{}
	crawler := &Crawler{parser: parser, domains: domains, mutex: mutex}

	for {
		select {
		case url := <-queue:
			// すでにキャッシュにある場合は何もしない
			if !cache.Contains(url) {
				continue
			}

			l.Wait(ctx)
			go crawler.Run(url)
		}
	}
}

/*
// GetLogger returns logger
func GetLogger() *zap.Logger {
	var once sync.Once
	var logger *zap.Logger
	once.Do(func() {
		logger, _ = zap.NewProduction()
	})
	return logger
}
*/

func (s *Spider) serve(urls []string) ([]*model.Ogp, []string) {
	ogps := make([]*model.Ogp, 0, len(urls))
	urlsRreturn := make([]string, 0, len(urls))

	for _, url := range urls {
		var ogp *model.Ogp

		ogp = s.serveCache(url)

		if ogp == nil {
			urlsRreturn = append(urlsRreturn, url)
		} else {
			ogps = append(ogps, ogp)
		}
	}

	return ogps, urlsRreturn
}

func (s *Spider) serveCache(url string) *model.Ogp {
	cached, ok := cache.Get(url)

	if ok {
		logURL("hit from cache", url)
		ogp := cached.(*model.Ogp)
		return ogp
	}

	logURL("cache miss", url)

	return nil
}

/*
func logError(msg string, tag string) {
	logger.Error(msg,
		zap.String("tag", tag),
	)
}

func logInfo(msg string, tag string) {
	logger.Info(msg,
		zap.String("tag", tag),
	)
}
*/

func logURL(msg string, url string) {
	fmt.Printf("%s: %s\n", msg, url)

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("Cannot open " + logFile + ": " + err.Error())
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags)

	log.Printf("%s: %s\n", msg, url)
	/*
		v := 0
		if requestToTop {
			v = 1
		}
		logger.Info(msg,
			zap.String("url", url),
			zap.String("tag", tag),
			zap.Int("request_to_top", v),
		)
	*/
}

/*
func logCount(msg string, count int, tag string) {
	logger.Info(msg,
		zap.Int("count", count),
		zap.String("tag", tag),
	)
}
*/
