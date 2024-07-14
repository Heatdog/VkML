package redis

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Heatdog/VkML/internal/models"
	"github.com/Heatdog/VkML/internal/storage"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
	logger *slog.Logger
}

func (s *Storage) Add(ctx context.Context, doc *models.Document) error {
	s.logger.Debug("add in cache", slog.Any("doc", doc))

	err := s.client.ZAdd(ctx, doc.URL, redis.Z{
		Member: doc,
		Score:  float64(doc.FetchTime),
	}).Err()

	if err != nil {
		s.logger.Warn(err.Error())
		return err
	}

	s.logger.Debug("insert in cache", slog.Any("doc", doc))
	return nil
}

func (s *Storage) GetByFetchTimeMax(ctx context.Context, url string) (models.Document, error) {
	s.logger.Debug("get doc by max fetch time", slog.String("url", url))

	scores, err := s.client.ZRevRangeWithScores(ctx, url, 0, 0).Result()
	if err != nil {
		s.logger.Warn(err.Error())
		return models.Document{}, err
	}

	if len(scores) < 1 {
		err = errors.New("no member selected")
		s.logger.Warn(err.Error())
		return models.Document{}, err
	}

	doc := scores[0].Member.(models.Document)
	s.logger.Debug("getted", slog.Any("doc", doc))

	return doc, nil
}

func (s *Storage) GetByFetchTimeMin(ctx context.Context, url string) (models.Document, error) {
	s.logger.Debug("get doc by min fetch time", slog.String("url", url))

	scores, err := s.client.ZRangeWithScores(ctx, url, 0, 0).Result()
	if err != nil {
		s.logger.Warn(err.Error())
		return models.Document{}, err
	}

	if len(scores) < 1 {
		err = errors.New("no member selected")
		s.logger.Warn(err.Error())
		return models.Document{}, err
	}

	doc := scores[0].Member.(models.Document)
	s.logger.Debug("getted", slog.Any("doc", doc))

	return doc, nil
}

func New(client *redis.Client, logger *slog.Logger) storage.Storage {
	return &Storage{
		client: client,
		logger: logger,
	}
}
