package service_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/Heatdog/VkML/internal/models"
	"github.com/Heatdog/VkML/internal/services"
	"github.com/Heatdog/VkML/internal/storage/postgre"
	"github.com/Heatdog/VkML/internal/storage/redis"
	"github.com/go-redis/redismock/v9"
	"github.com/pashagolub/pgxmock/v3"
	lib_redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer dbMock.Close()

	cacheDb, mock := redismock.NewClientMock()
	defer cacheDb.Close()

	opt := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelError,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opt))
	slog.SetDefault(logger)

	storage := postgre.New(dbMock, logger)
	cache := redis.New(cacheDb, logger)
	documentProcess := services.New(storage, cache, logger)

	type mockBehavior func(inDocument models.Document, outDocument *models.Document, err error)

	testTable := []struct {
		name        string
		mockFunc    mockBehavior
		inDocument  models.Document
		outDocument *models.Document
		err         error
	}{
		{
			name: "ok get from database",
			inDocument: models.Document{
				URL:       "/set",
				PubDate:   12345,
				FetchTime: 45567,
				Text:      "1234",
			},
			outDocument: &models.Document{
				URL:            "/set",
				PubDate:        12300,
				FetchTime:      99999,
				Text:           "5668",
				FirstFetchTime: 12300,
			},
			err: nil,
			mockFunc: func(inDocument models.Document, outDocument *models.Document, err error) {
				dbMock.ExpectExec("INSERT INTO documents").WithArgs(inDocument.URL,
					inDocument.PubDate, inDocument.FetchTime, inDocument.Text).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mock.ExpectZAdd(inDocument.URL, lib_redis.Z{
					Member: &inDocument,
					Score:  float64(inDocument.FetchTime),
				}).SetVal(1)

				mock.ExpectZRangeWithScores(inDocument.URL, 0, 0).
					SetVal([]lib_redis.Z{})

				rowMin := pgxmock.NewRows([]string{"pub_date", "fetch_time", "text"})
				rowMin.AddRow(outDocument.PubDate, outDocument.FirstFetchTime,
					"hello world")

				dbMock.ExpectQuery(`
					SELECT pub_date, fetch_time, text 
					FROM documents 
					WHERE url`).
					WithArgs(inDocument.URL).WillReturnRows(rowMin)

				mock.ExpectZRevRangeWithScores(inDocument.URL, 0, 0).
					SetVal([]lib_redis.Z{})

				rowMax := pgxmock.NewRows([]string{"pub_date", "fetch_time", "text"})
				rowMax.AddRow(uint64(45544), outDocument.FetchTime, outDocument.Text)

				dbMock.ExpectQuery(`
					SELECT pub_date, fetch_time, text
					FROM documents
					WHERE url`).
					WithArgs(inDocument.URL).WillReturnRows(rowMax)

			},
		},
		{
			name: "ok get from cache",
			inDocument: models.Document{
				URL:       "/set",
				PubDate:   12345,
				FetchTime: 45567,
				Text:      "1234",
			},
			outDocument: &models.Document{
				URL:            "/set",
				PubDate:        12300,
				FetchTime:      99999,
				Text:           "5668",
				FirstFetchTime: 12300,
			},
			err: nil,
			mockFunc: func(inDocument models.Document, outDocument *models.Document, err error) {
				dbMock.ExpectExec("INSERT INTO documents").WithArgs(inDocument.URL,
					inDocument.PubDate, inDocument.FetchTime, inDocument.Text).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mock.ExpectZAdd(inDocument.URL, lib_redis.Z{
					Member: &inDocument,
					Score:  float64(inDocument.FetchTime),
				}).SetVal(1)

				mock.ExpectZRangeWithScores(inDocument.URL, 0, 0).
					SetVal([]lib_redis.Z{
						{
							Member: models.Document{
								PubDate:   outDocument.PubDate,
								FetchTime: outDocument.FirstFetchTime,
							},
						},
					})

				mock.ExpectZRevRangeWithScores(inDocument.URL, 0, 0).
					SetVal([]lib_redis.Z{
						{
							Member: *outDocument,
						},
					})
			},
		},
		{
			name: "insert error",
			inDocument: models.Document{
				URL:       "/set",
				PubDate:   12345,
				FetchTime: 45567,
				Text:      "1234",
			},
			outDocument: nil,
			err:         errors.New("insert error"),
			mockFunc: func(inDocument models.Document, outDocument *models.Document, err error) {
				dbMock.ExpectExec("INSERT INTO documents").WithArgs(inDocument.URL,
					inDocument.PubDate, inDocument.FetchTime, inDocument.Text).
					WillReturnError(err)
			},
		},
		{
			name: "insert affected zero rows",
			inDocument: models.Document{
				URL:       "/set",
				PubDate:   12345,
				FetchTime: 45567,
				Text:      "1234",
			},
			outDocument: nil,
			err:         errors.New("zero rows affected"),
			mockFunc: func(inDocument models.Document, outDocument *models.Document, err error) {
				dbMock.ExpectExec("INSERT INTO documents").WithArgs(inDocument.URL,
					inDocument.PubDate, inDocument.FetchTime, inDocument.Text).
					WillReturnResult(pgxmock.NewResult("INSERT", 0))
			},
		},
		{
			name: "select min error",
			inDocument: models.Document{
				URL:       "/set",
				PubDate:   12345,
				FetchTime: 45567,
				Text:      "1234",
			},
			outDocument: nil,
			err:         errors.New("select error"),
			mockFunc: func(inDocument models.Document, outDocument *models.Document, err error) {
				dbMock.ExpectExec("INSERT INTO documents").WithArgs(inDocument.URL,
					inDocument.PubDate, inDocument.FetchTime, inDocument.Text).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mock.ExpectZAdd(inDocument.URL, lib_redis.Z{
					Member: &inDocument,
					Score:  float64(inDocument.FetchTime),
				}).SetVal(1)

				dbMock.ExpectQuery(`
					SELECT pub_date, fetch_time, text 
					FROM documents 
					WHERE url`).
					WithArgs(inDocument.URL).WillReturnError(err)
			},
		},
		{
			name: "select max error",
			inDocument: models.Document{
				URL:       "/set",
				PubDate:   12345,
				FetchTime: 45567,
				Text:      "1234",
			},
			outDocument: nil,
			err:         errors.New("select error"),
			mockFunc: func(inDocument models.Document, outDocument *models.Document, err error) {
				dbMock.ExpectExec("INSERT INTO documents").WithArgs(inDocument.URL,
					inDocument.PubDate, inDocument.FetchTime, inDocument.Text).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mock.ExpectZAdd(inDocument.URL, lib_redis.Z{
					Member: &inDocument,
					Score:  float64(inDocument.FetchTime),
				}).SetVal(1)

				rowMin := pgxmock.NewRows([]string{"pub_date", "fetch_time", "text"})
				rowMin.AddRow(uint64(123), uint64(123), "hello world")

				dbMock.ExpectQuery(`
						SELECT pub_date, fetch_time, text 
						FROM documents 
						WHERE url`).
					WithArgs(inDocument.URL).WillReturnRows(rowMin)

				dbMock.ExpectQuery(`
						SELECT pub_date, fetch_time, text
						FROM documents
						WHERE url`).
					WithArgs(inDocument.URL).WillReturnError(err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(testCase.inDocument, testCase.outDocument, testCase.err)
			res, err := documentProcess.Process(&testCase.inDocument)

			require.Equal(t, testCase.err, err)
			require.Equal(t, testCase.outDocument, res)
		})
	}
}
