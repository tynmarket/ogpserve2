package spider

import (
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	"github.com/tynmarket/ogpserve2/model"
	"golang.org/x/net/html/charset"
)

// Parser is
type Parser struct {
}

// OGP meta tags
const (
	TYPE        = "type"
	TITLE       = "title"
	DESCRIPTION = "description"
	URL         = "url"
	IMAGE       = "image"
)

// Twitter card meta tags
const (
	CARD    = "card"
	SITE    = "site"
	PHOTO   = "photo"
	SUMMARY = "summary"
	LARGE   = "summary_large_image"
	DEFAULT = "https://tyn-imarket.com/public/default-image.png"
)

var reCharset = regexp.MustCompile("meta charset=\"(.+)\"")
var reCharsetOld = regexp.MustCompile("meta.+content=\".*charset=(.+)\"")
var reOgp = regexp.MustCompile("(?s)meta property=\"og:([a-z]+)\"[\n ]+content=[\"']([^<>]+?)[\"']")
var reOgpRev = regexp.MustCompile("(?s)meta content=[\"']([^<>]+?)[\"'][\n ]+property=\"og:([a-z]+)\"")
var reCard = regexp.MustCompile("meta (name|property)=\"twitter:([a-z]+)\" content=\"([^<>]+?)\"")
var reCardRev = regexp.MustCompile("meta content=\"([^<>]+?)\" (name|property)=\"twitter:([a-z]+)\"")
var reTitleTag = regexp.MustCompile("<title.*>(.+)</title>")
var reFigureTag = regexp.MustCompile("<figure.*><img .*src=\"(.+?)\" .*></figure>")

// Parse ogp meta tags
func (p *Parser) parse(requestURL string, html string) {
	html = convert(html)
	ogps := reOgp.FindAllStringSubmatch(html, -1)
	ogpsRev := reOgpRev.FindAllStringSubmatch(html, -1)
	cards := reCard.FindAllStringSubmatch(html, -1)
	cardsRev := reCardRev.FindAllStringSubmatch(html, -1)

	ogp := &model.Ogp{TwitterCard: &model.TwitterCard{}}

	// ogp
	for _, v := range ogps {
		switch v[1] {
		case TYPE:
			ogp.Type = v[2]
		case TITLE:
			ogp.Title = v[2]
		case DESCRIPTION:
			ogp.Description = v[2]
		case URL:
			ogp.URL = v[2]
		case IMAGE:
			ogp.Image = v[2]
		}
	}
	// ogp reverse order
	for _, v := range ogpsRev {
		switch v[2] {
		case TYPE:
			ogp.Type = v[1]
		case TITLE:
			ogp.Title = v[1]
		case DESCRIPTION:
			ogp.Description = v[1]
		case URL:
			ogp.URL = v[1]
		case IMAGE:
			ogp.Image = v[1]
		}
	}

	// Twitter Card
	for _, v := range cards {
		switch v[2] {
		case CARD:
			ogp.TwitterCard.Card = v[3]
		case SITE:
			ogp.TwitterCard.Site = v[3]
		case TITLE:
			ogp.TwitterCard.Title = v[3]
		case DESCRIPTION:
			ogp.TwitterCard.Description = v[3]
		case IMAGE:
			ogp.TwitterCard.Image = v[3]
		}
	}

	// Twitter Card reverse order
	for _, v := range cardsRev {
		switch v[3] {
		case CARD:
			ogp.TwitterCard.Card = v[1]
		case SITE:
			ogp.TwitterCard.Site = v[1]
		case TITLE:
			ogp.TwitterCard.Title = v[1]
		case DESCRIPTION:
			ogp.TwitterCard.Description = v[1]
		case IMAGE:
			ogp.TwitterCard.Image = v[1]
		}
	}

	// Parse title from title tag when no title meta tags found
	if ogp.Title == "" && ogp.TwitterCard.Title == "" {
		match := reTitleTag.FindStringSubmatch(html)

		if len(match) > 1 {
			ogp.Title = match[1]
		}
	}

	// figureをimageとして取得
	if ogp.Image == "" && ogp.TwitterCard.Image == "" {
		match := reFigureTag.FindStringSubmatch(html)

		if len(match) > 1 {
			ogp.Image = match[1]
			ogp.TwitterCard.Card = LARGE
		}
	}

	// Convert relative path to absolute path
	if strings.HasPrefix(ogp.Image, "/") {
		url, err := url.Parse(ogp.URL)
		if err == nil {
			ogp.Image = url.Scheme + "://" + url.Host + ogp.Image
		}
	}
	if strings.HasPrefix(ogp.TwitterCard.Image, "/") {
		url, err := url.Parse(ogp.URL)
		if err == nil {
			ogp.TwitterCard.Image = url.Scheme + "://" + url.Host + ogp.TwitterCard.Image
		}
	}

	if ogp.TwitterCard.Card == PHOTO ||
		ogp.Title != "" && ogp.TwitterCard.Card == "" {
		ogp.TwitterCard.Card = SUMMARY
	}

	if (ogp.Title != "" || ogp.TwitterCard.Title != "") &&
		ogp.Image == "" && ogp.TwitterCard.Image == "" {
		ogp.Image = DEFAULT
		ogp.TwitterCard.Card = SUMMARY
	}

	ogp.RequestURL = requestURL

	cache.Add(requestURL, ogp)
}

func convert(html string) string {
	match := reCharset.FindStringSubmatch(html)
	if len(match) == 0 {
		match = reCharsetOld.FindStringSubmatch(html)

		if len(match) == 0 {
			return html
		}
	}

	char := match[1]
	_, name := charset.Lookup(char)

	if name == "shift_jis" || name == "euc-jp" {
		return convertEncoding(html, char)
	}

	return html
}

func convertEncoding(title string, name string) string {
	s := strings.NewReader(title)
	reader, err := charset.NewReaderLabel(name, s)
	if err != nil {
		return title
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return title
	}
	return string(bytes)
}
