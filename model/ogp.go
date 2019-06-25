package model

// Ogp represents Opn Graph Protocol including Twitter Card
type Ogp struct {
	Type        string       `json:"type"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	URL         string       `json:"url"`
	Image       string       `json:"image"`
	Retry       bool         `json:"retry"`
	TwitterCard *TwitterCard `json:"twitter"`
	SummaryCard bool
}

// Copy returns shallow copy of Ogp
func (o *Ogp) Copy() *Ogp {
	return &Ogp{
		Type:        o.Type,
		Title:       o.Title,
		Description: o.Description,
		URL:         o.URL,
		Image:       o.Image,
		TwitterCard: o.TwitterCard,
	}
}

// MergeIntoTwitter merge OGP meta tags into TwitterCard
func (o *Ogp) MergeIntoTwitter() *TwitterCard {
	if o.TwitterCard.Title == "" {
		o.TwitterCard.Title = o.Title
	}
	if o.TwitterCard.Description == "" {
		o.TwitterCard.Description = o.Description
	}
	if o.TwitterCard.Image == "" {
		o.TwitterCard.Image = o.Image
	}
	if o.SummaryCard || (o.TwitterCard.Card == "" && o.TwitterCard.Image != "") {
		o.TwitterCard.Card = summaryCard
	}
	if o.Retry {
		o.TwitterCard.Retry = o.Retry
	}
	return o.TwitterCard
}
