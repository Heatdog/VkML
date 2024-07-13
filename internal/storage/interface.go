package storage

import (
	"context"

	"github.com/Heatdog/VkML/internal/models"
)

type Storage interface {
	Add(ctx context.Context, doc *models.Document) error
	GetByFetchTimeMin(ctx context.Context, url string) (models.Document, error)
	GetByFetchTimeMax(ctx context.Context, url string) (models.Document, error)
}
