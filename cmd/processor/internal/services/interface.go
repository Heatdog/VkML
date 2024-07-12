package services

import "github.com/Heatdog/VkML/cmd/processor/internal/models"

type Processor interface {
	Process(d *models.Document) (*models.Document, error)
}
