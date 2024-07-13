package postgre

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Heatdog/VkML/internal/models"
	"github.com/Heatdog/VkML/internal/storage"
	"github.com/Heatdog/VkML/pkg/storage/postgre"
)

type Storage struct {
	client postgre.Client
	logger *slog.Logger
}

func New(client postgre.Client, logger *slog.Logger) storage.Storage {
	return &Storage{
		client: client,
		logger: logger,
	}
}

func (storage *Storage) Add(ctx context.Context, doc *models.Document) error {
	storage.logger.Debug("insert document")

	q := `
		INSERT INTO documents (url, pub_date, fetch_time, text)
		VALUES ($1, $2, $3, $4)
	`

	storage.logger.Debug("document", slog.Any("doc", doc))
	tag, err := storage.client.Exec(ctx, q, doc.URL, doc.PubDate, doc.FetchTime, doc.Text)
	if err != nil {
		storage.logger.Warn("error", err.Error())
		return err
	}

	if tag.RowsAffected() == 0 {
		storage.logger.Debug("zero rows affected")
		return errors.New("zero rows affected")
	}

	storage.logger.Debug("successful document insert")
	return nil
}

func (storage *Storage) GetByFetchTimeMin(ctx context.Context, url string) (models.Document, error) {
	storage.logger.Debug("get min fetch time document")

	q := `
		SELECT pub_date, fetch_time, text
		FROM documents
		WHERE url = $1
		ORDER BY fetch_time ASC
		LIMIT 1
	`

	storage.logger.Debug("query", q)
	row := storage.client.QueryRow(ctx, q, url)

	var document models.Document
	if err := row.Scan(&document.PubDate, &document.FetchTime, &document.Text); err != nil {
		storage.logger.Warn(err.Error())
		return document, err
	}

	document.URL = url

	storage.logger.Debug("selected document", slog.Any("doc", document))
	return document, nil
}

func (storage *Storage) GetByFetchTimeMax(ctx context.Context, url string) (models.Document, error) {
	storage.logger.Debug("get max fetch time document")

	q := `
		SELECT pub_date, fetch_time, text
		FROM documents
		WHERE url = $1
		ORDER BY fetch_time DESC
		LIMIT 1
	`

	storage.logger.Debug("query", q)
	row := storage.client.QueryRow(ctx, q, url)

	var document models.Document
	if err := row.Scan(&document.PubDate, &document.FetchTime, &document.Text); err != nil {
		storage.logger.Warn(err.Error())
		return document, err
	}

	document.URL = url

	storage.logger.Debug("selected document", slog.Any("doc", document))
	return document, nil
}
