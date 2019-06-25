package spider

import (
	"strings"
	"sync"
)

// Crawler struct
type Crawler struct {
	domains map[string]struct{}
	mutex   *sync.Mutex
}

// Run is
func (c *Crawler) Run(url string) {
	ok := c.lockDomain(url)

	if ok {
		logURL("crawl", url)

		c.freeDomain(url)
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
		_, ok := c.domains[domain]

		if !ok {
			c.domains[domain] = struct{}{}
			return true
		}
	}

	return false
}

func (c *Crawler) freeDomain(url string) {
	domain := getDomain(url)

	if domain != "" {
		defer c.mutex.Unlock()
		c.mutex.Lock()
		delete(c.domains, domain)
	}
}

func getDomain(url string) string {
	strs := strings.Split(url, "/")
	if len(strs) > 2 {
		return strs[2]
	}
	return ""
}
