package models

type CustomEmbed struct {
	URL       string
	Title     string
	Desc      string
	Type      string
	Timestamp string
	Color     int

	FooterText     string
	FooterImageURL string

	ImageURL string
	ImageH   int
	ImageW   int

	ThumbnailURL string
	ThumbnailH   int
	ThumbnailW   int

	ProviderURL  string
	ProviderText string

	AuthorText     string
	AuthorURL      string
	AuthorImageURL string
}
