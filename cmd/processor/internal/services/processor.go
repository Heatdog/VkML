package services

import (
	"log/slog"

	"github.com/Heatdog/VkML/cmd/processor/internal/models"
	"github.com/Heatdog/VkML/cmd/processor/internal/storage"
)

type ProcessorDocuments struct {
	logger  *slog.Logger
	storage storage.Storage
}

func New(storage storage.Storage) Processor {
	return &ProcessorDocuments{
		storage: storage,
	}
}

func (processor *ProcessorDocuments) Process(d *models.Document) (*models.Document, error) {

	return nil, nil
}
