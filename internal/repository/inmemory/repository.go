package inmemory

import (
	"context"
	"main/internal/pkg/errs"
	"main/internal/repository/entities"
	"sync"
)

type UrlRepository struct {
	shortToOrigin sync.Map
	originToShort sync.Map
	mu            sync.Mutex
}

func New() *UrlRepository {
	return &UrlRepository{}
}

func (r *UrlRepository) AddUrlRelation(ctx context.Context, req entities.AddUrlRelationRequest) (entities.AddUrlRelationResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.shortToOrigin.Load(req.ShortenedUrl); exist {
		return entities.AddUrlRelationResponse{}, errs.ErrorRepositoryDuplicate
	}

	if _, exist := r.originToShort.Load(req.OriginUrl); exist {
		return entities.AddUrlRelationResponse{}, errs.ErrorAlreadyExist
	}

	r.shortToOrigin.Store(req.ShortenedUrl, req.OriginUrl)
	r.originToShort.Store(req.OriginUrl, req.ShortenedUrl)

	return entities.AddUrlRelationResponse{ShortenedUrl: req.ShortenedUrl}, nil
}

func (r *UrlRepository) GetOriginURLFromShortened(ctx context.Context, req entities.GetOriginURLFromShortenedUrlRequest) (entities.GetOriginURLFromShortenedUrlResponse, error) {
	origin, exist := r.shortToOrigin.Load(req.ShortenedUrl)
	if !exist {
		return entities.GetOriginURLFromShortenedUrlResponse{OriginUrl: ""}, errs.ErrorRepositoryUrlEmpty
	}

	val, _ := origin.(string)
	return entities.GetOriginURLFromShortenedUrlResponse{OriginUrl: val}, nil
}

func (r *UrlRepository) Close() error {
	return nil
}
