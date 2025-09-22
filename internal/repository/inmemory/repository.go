package inmemory

import (
	"context"
	"main/internal/errs"
	"main/internal/repository/entities"
)

type UrlRepository struct {
	dict map[string]string
}

func New() *UrlRepository {
	return &UrlRepository{dict: make(map[string]string)}
}

func (r *UrlRepository) AddUrlRelation(ctx context.Context, req entities.AddUrlRelationRequest) (entities.AddUrlRelationResponse, error) {
	if _, exist := r.dict[req.ShortenedUrl]; exist {
		return entities.AddUrlRelationResponse{}, errs.ErrorRepositoryDuplicate
	}
	for _, v := range r.dict {
		if v == req.OriginUrl {
			return entities.AddUrlRelationResponse{}, errs.ErrorAlreadyExist
		}
	}
	r.dict[req.ShortenedUrl] = req.OriginUrl
	return entities.AddUrlRelationResponse{ShortenedUrl: req.ShortenedUrl}, nil
}

func (r *UrlRepository) GetOriginURLFromShortened(ctx context.Context, req entities.GetOriginURLFromShortenedUrlRequest) (entities.GetOriginURLFromShortenedUrlResponse, error) {
	origin, exist := r.dict[req.ShortenedUrl]
	if !exist {
		return entities.GetOriginURLFromShortenedUrlResponse{""}, errs.ErrorRepositoryUrlEmpty
	}
	return entities.GetOriginURLFromShortenedUrlResponse{origin}, nil
}

func (r *UrlRepository) Close() error {
	return nil
}
