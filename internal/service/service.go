package service

import (
	"context"
	"errors"
	"main/internal/pkg/errs"
	"main/internal/pkg/logger"
	"main/internal/repository/entities"
	serviceEntities "main/internal/service/entities"
	"main/pkg/cache"
)

type UrlRepository interface {
	AddUrlRelation(ctx context.Context, req entities.AddUrlRelationRequest) (entities.AddUrlRelationResponse, error)
	GetOriginURLFromShortened(ctx context.Context, req entities.GetOriginURLFromShortenedUrlRequest) (entities.GetOriginURLFromShortenedUrlResponse, error)
	Close() error
}

type Generator interface {
	Generate() string
	BaseHost() string
}
type UrlService struct {
	generator Generator
	repo      UrlRepository
	cach      cache.Cache[string, string]
	log       logger.Logger
}

func NewUrlService(urlGen Generator, repo UrlRepository, cach cache.Cache[string, string], log logger.Logger) *UrlService {
	return &UrlService{
		repo:      repo,
		generator: urlGen,
		cach:      cach,
		log:       log,
	}
}

func (s *UrlService) GenerateNewLink(ctx context.Context, req serviceEntities.GenerateLinkRequest) (serviceEntities.GenerateLinkResponse, error) {
	s.log.Debugw("GenerateNewLink", "ctx", ctx, "req", req)

	generatedUrl, err := s.repo.AddUrlRelation(
		ctx,
		entities.AddUrlRelationRequest{
			OriginUrl:    req.OriginUrl,
			ShortenedUrl: s.generator.Generate(),
		})
	if errors.Is(err, errs.ErrorAlreadyExist) {
		return serviceEntities.GenerateLinkResponse{}, errs.ErrorAlreadyExist
	}

	for errors.Is(err, errs.ErrorRepositoryDuplicate) {
		s.log.Errorw("GenerateNewLink", "err", err)
		generatedUrl, err = s.repo.AddUrlRelation(
			ctx,
			entities.AddUrlRelationRequest{
				OriginUrl:    req.OriginUrl,
				ShortenedUrl: s.generator.Generate(),
			})
	}
	if err != nil {
		s.log.Errorw("GenerateNewLink", "err", err)
		return serviceEntities.GenerateLinkResponse{}, errs.ErrorUrlServiceInternal
	}

	s.log.Debugw("GenerateNewLink", "generatedUrl", generatedUrl)

	s.cach.Set(generatedUrl.ShortenedUrl, req.OriginUrl)

	return serviceEntities.GenerateLinkResponse{ShortenedUrl: s.generator.BaseHost() + generatedUrl.ShortenedUrl}, nil
}

func (s *UrlService) GetLink(ctx context.Context, req serviceEntities.GetLinkRequest) (serviceEntities.GetLinkResponse, error) {
	s.log.Debugw("GetLink start", "ctx", ctx, "short url", req.ShortenedUrl)

	origin, ok := s.cach.Get(req.ShortenedUrl)
	if ok {
		return serviceEntities.GetLinkResponse{OriginUrl: origin}, nil
	}

	link, err := s.repo.GetOriginURLFromShortened(ctx, entities.GetOriginURLFromShortenedUrlRequest{ShortenedUrl: req.ShortenedUrl})

	if err != nil {
		s.log.Errorw("GetLink", "err", err)

		if errors.Is(err, errs.ErrorRepositoryUrlEmpty) {
			return serviceEntities.GetLinkResponse{}, errs.ErrorUrlServiceLinkNotFound
		} else {
			return serviceEntities.GetLinkResponse{}, errs.ErrorUrlServiceInternal
		}
	}
	s.log.Debugw("GetLink end", "ctx", ctx, "origin url", link)

	s.cach.Set(req.ShortenedUrl, link.OriginUrl)

	return serviceEntities.GetLinkResponse{OriginUrl: link.OriginUrl}, nil
}
