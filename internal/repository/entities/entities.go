package entities

type AddUrlRelationRequest struct {
	OriginUrl    string
	ShortenedUrl string
}
type AddUrlRelationResponse struct {
	ShortenedUrl string
}

type GetOriginURLFromShortenedUrlRequest struct {
	ShortenedUrl string
}
type GetOriginURLFromShortenedUrlResponse struct {
	OriginUrl string
}
