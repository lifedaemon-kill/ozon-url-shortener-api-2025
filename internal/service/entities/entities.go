package entities

type GenerateLinkRequest struct {
	OriginUrl string `json:"url"`
}
type GenerateLinkResponse struct {
	ShortenedUrl string `json:"url"`
}

type GetLinkRequest struct {
	ShortenedUrl string `json:"url"`
}
type GetLinkResponse struct {
	OriginUrl string `json:"url"`
}
