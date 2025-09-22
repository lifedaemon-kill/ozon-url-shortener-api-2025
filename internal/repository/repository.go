package repository

import (
	"context"
	"main/internal/repository/entities"
)

type UrlRepository interface {
	AddUrlRelation(ctx context.Context, req entities.AddUrlRelationRequest) (entities.AddUrlRelationResponse, error)
	GetOriginURLFromShortened(ctx context.Context, req entities.GetOriginURLFromShortenedUrlRequest) (entities.GetOriginURLFromShortenedUrlResponse, error)
	Close() error
}
