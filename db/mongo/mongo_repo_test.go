package mongo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/go-ota-svc-common/db/mongo"
)

func TestMongoDB(t *testing.T) {
	t.Run("should get ErrorConnection if connection string is Invalid", func(t *testing.T) {
		var (
			connStr = "INVALID_CONNECTION_STRING"
		)
		ctx := context.Background()
		log := logger.NewNopLogger()
		_, err := mongo.NewMongoDB(ctx, log, connStr)
		var typedErr apperrors.AppError
		if err == nil || errors.As(err, &typedErr) && typedErr.ErrorCode != apperrors.ErrorDbConnection {
			t.Errorf("got %s, expected %s", err, apperrors.ErrorDbConnection)
		}
	})
}
