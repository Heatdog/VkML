package service_test

import (
	"log/slog"
	"os"
	"strconv"
	"testing"

	"github.com/Heatdog/VkML/internal/models"
	"github.com/Heatdog/VkML/internal/services"
	"github.com/Heatdog/VkML/internal/storage/postgre"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer dbMock.Close()

	opt := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelError,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opt))
	slog.SetDefault(logger)

	storage := postgre.New(dbMock, logger)
	documentProcess := services.New(storage, logger)

	type mockBehavior func(inDocument models.Document, outDocument models.Document)

	testTable := []struct {
		name        string
		mockFunc    mockBehavior
		inDocument  models.Document
		outDocument models.Document
		err         error
	}{
		{
			name: "ok",
			inDocument: models.Document{
				URL:       "/set",
				PubDate:   12345,
				FetchTime: 45567,
				Text:      "1234",
			},
			outDocument: models.Document{
				URL:            "/set",
				PubDate:        12300,
				FetchTime:      99999,
				Text:           "5668",
				FirstFetchTime: 12300,
			},
			err: nil,
			mockFunc: func(inDocument models.Document, outDocument models.Document) {
				dbMock.ExpectExec("INSERT INTO documents").WithArgs(inDocument.URL,
					inDocument.PubDate, inDocument.FetchTime, inDocument.Text).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				rowMin := pgxmock.NewRows([]string{"pub_date", "fetch_time", "text"})
				rowMin.AddRow(strconv.Itoa(int(outDocument.PubDate)), strconv.Itoa(int(outDocument.FirstFetchTime)),
					"hello world")
				rowMax := pgxmock.NewRows([]string{"pub_date", "fetch_time", "text"})
				rowMax.AddRow("45544", strconv.Itoa(int(outDocument.FetchTime)), outDocument.Text)

				dbMock.ExpectQuery(`
					SELECT pub_date, fetch_time, text
					FROM documents
					WHERE url
					`).WithArgs(inDocument.URL).WillReturnRows(rowMin)

				dbMock.ExpectQuery(`
					SELECT pub_date, fetch_time, text
					FROM documents
					WHERE url
					`).WithArgs(inDocument.URL).WillReturnRows(rowMax)

			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(testCase.inDocument, testCase.outDocument)
			res, err := documentProcess.Process(&testCase.inDocument)

			require.Equal(t, testCase.err, err)
			require.Equal(t, testCase.outDocument, *res)
		})
	}
}
