package grpc

import (
	"context"
	"errors"
	"main/internal/pkg/errs"
	"main/internal/pkg/logger"
	"main/internal/service/entities"
	serviceEntities "main/internal/service/entities"
	desc "main/pkg/protogen/url-shortener"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	GenerateNewLink(context.Context, serviceEntities.GenerateLinkRequest) (serviceEntities.GenerateLinkResponse, error)
	GetLink(context.Context, serviceEntities.GetLinkRequest) (serviceEntities.GetLinkResponse, error)
}

type Implementation struct {
	desc.UnimplementedUrlShortenerServer

	urlService Service
	log        logger.Logger
}

func New(urlService Service, log logger.Logger) *Implementation {
	return &Implementation{urlService: urlService, log: log}
}

func (i *Implementation) WrapLink(ctx context.Context, req *desc.OriginURL) (*desc.ShortURL, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	res, err := i.urlService.GenerateNewLink(ctx, entities.GenerateLinkRequest{req.Url})
	if err != nil {
		i.log.Errorw("Error generating new link", "url", req.Url, "err", err)

		if errors.Is(err, errs.ErrorAlreadyExist) {
			return nil, errs.ErrorAlreadyExist
		}
		return nil, status.Errorf(codes.Internal, "Error generating new link")
	}
	return &desc.ShortURL{Url: res.ShortenedUrl}, nil
}

func (i *Implementation) UnwrapLink(ctx context.Context, req *desc.ShortURL) (*desc.OriginURL, error) {
	res, err := i.urlService.GetLink(ctx, entities.GetLinkRequest{req.Url})
	if err != nil {
		i.log.Errorw("Error getting link", "url", req.Url, "err", err)
		if errors.Is(err, errs.ErrorUrlServiceLinkNotFound) {
			return nil, errs.ErrorUrlServiceLinkNotFound
		}
		return nil, status.Errorf(codes.Internal, "Error getting link")
	}
	return &desc.OriginURL{Url: res.OriginUrl}, nil
}
