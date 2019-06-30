package spider

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Crawler struct
type Crawler struct {
	parser  *Parser
	domains map[string]int64
	mutex   *sync.Mutex
}

const (
	intervalMills = int64(1000) // 1000 mills
)

// Run is
func (c *Crawler) Run(url string) {
	ok := c.lockDomain(url)

	if ok {
		logURL("crawl for", url)

		c.crawl(url)
	} else {
		logURL("add to queue", url)
		queue <- url
	}

}

func (c *Crawler) lockDomain(url string) bool {
	defer c.mutex.Unlock()
	domain := getDomain(url)

	if domain != "" {
		c.mutex.Lock()

		// キャッシュが溜まったら常にクリア
		c.checkCacheSizeAndClear()

		prevNextTime, ok := c.domains[domain]
		current := time.Now().UnixNano() / 1000000

		// 初回またはインターバル経過後
		if !ok || current >= prevNextTime {
			nextTime := current + intervalMills
			c.domains[domain] = nextTime
			return true
		}
	}

	return false
}

func (c *Crawler) checkCacheSizeAndClear() {
	if cache.Len() >= cacheSize {
		cache.Purge()
		c.domains = make(map[string]int64)
	}
}

func getDomain(url string) string {
	strs := strings.Split(url, "/")
	if len(strs) > 2 {
		return strs[2]
	}
	return ""
}

func (c *Crawler) crawl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	html := string(bytes)

	// Send to Parser
	c.parser.parse(url, html)
}
