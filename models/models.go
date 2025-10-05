package models

type AliasShortener struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

type UrlShortener struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}
