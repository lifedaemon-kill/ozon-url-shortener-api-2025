package inmemory

import (
	"context"
	"main/internal/errs"
	"main/internal/repository/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryUrlRepository_AddUrlRelation(t *testing.T) {
	ctx := context.Background()
	repo := New()

	t.Run("успешное добавление новой ссылки", func(t *testing.T) {
		resp, err := repo.AddUrlRelation(ctx, entities.AddUrlRelationRequest{
			OriginUrl:    "https://example.com",
			ShortenedUrl: "abc123",
		})
		assert.NoError(t, err)
		assert.Equal(t, "abc123", resp.ShortenedUrl)
	})

	t.Run("дубликат сокращённой ссылки", func(t *testing.T) {
		_, err := repo.AddUrlRelation(ctx, entities.AddUrlRelationRequest{
			OriginUrl:    "https://another.com",
			ShortenedUrl: "abc123",
		})
		assert.ErrorIs(t, err, errs.ErrorRepositoryDuplicate)
	})

	t.Run("дубликат оригинальной ссылки", func(t *testing.T) {
		_, err := repo.AddUrlRelation(ctx, entities.AddUrlRelationRequest{
			OriginUrl:    "https://example.com",
			ShortenedUrl: "newshort",
		})
		assert.ErrorIs(t, err, errs.ErrorAlreadyExist)
	})
}

func TestInMemoryUrlRepository_GetOriginURLFromShortened(t *testing.T) {
	ctx := context.Background()
	repo := New()

	_, _ = repo.AddUrlRelation(ctx, entities.AddUrlRelationRequest{
		OriginUrl:    "https://example.com",
		ShortenedUrl: "abc123",
	})

	t.Run("успешное получение оригинального URL", func(t *testing.T) {
		resp, err := repo.GetOriginURLFromShortened(ctx, entities.GetOriginURLFromShortenedUrlRequest{
			ShortenedUrl: "abc123",
		})
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", resp.OriginUrl)
	})

	t.Run("не найдено URL", func(t *testing.T) {
		_, err := repo.GetOriginURLFromShortened(ctx, entities.GetOriginURLFromShortenedUrlRequest{
			ShortenedUrl: "notexist",
		})
		assert.ErrorIs(t, err, errs.ErrorRepositoryUrlEmpty)
	})
}
