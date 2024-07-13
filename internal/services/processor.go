package services

import (
	"context"
	"log/slog"

	"github.com/Heatdog/VkML/internal/models"
	"github.com/Heatdog/VkML/internal/storage"
)

type DocumentsProcessor struct {
	logger  *slog.Logger
	storage storage.Storage
}

func New(storage storage.Storage, logger *slog.Logger) Processor {
	return &DocumentsProcessor{
		storage: storage,
		logger:  logger,
	}
}

func (processor *DocumentsProcessor) Process(d *models.Document) (*models.Document, error) {
	processor.logger.Debug("document", slog.Any("params", d))

	ctx := context.TODO() // по-хорошему передавать контекст в саму функцию Process извне

	if err := processor.storage.Add(ctx, d); err != nil {
		processor.logger.Warn(err.Error())
		return nil, err
	}

	minDoc, err := processor.storage.GetByFetchTimeMin(ctx, d.URL)
	if err != nil {
		processor.logger.Warn(err.Error())
		return nil, err
	}

	d.PubDate = minDoc.PubDate
	d.FirstFetchTime = minDoc.FetchTime

	maxDoc, err := processor.storage.GetByFetchTimeMax(ctx, d.URL)
	if err != nil {
		processor.logger.Warn(err.Error())
		return nil, err
	}

	d.Text = maxDoc.Text
	d.FetchTime = maxDoc.FetchTime

	processor.logger.Debug("result document", slog.Any("doc", d))

	return d, nil
}
