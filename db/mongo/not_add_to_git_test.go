package mongo_test

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/db/mongo"
)

type testDocument struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Version int                `bson:"version"`
	Data    json.RawMessage    `bson:"data"`
}

func TestMongoDB1(t *testing.T) {
	t.Run("should get ErrorConnection if connection string is Invalid", func(t *testing.T) {
		var connStr = "mongodb://mongoadmin:secret@localhost:27017/test?authSource=admin"
		ctx := context.Background()
		log := logger.NewNopLogger()
		mdb, err := mongo.NewMongoDB(ctx, log, connStr)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		coll1 := mdb.GetCollection("test_coll11")
		data := "this is test string blablablablablablabla"
		d := testDocument{
			ID:      primitive.NewObjectID(),
			Version: 10,
			Data:    json.RawMessage(data),
		}
		id, err := mdb.InsertOne(ctx, coll1, d)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		var d2 testDocument
		if err = mdb.GetOneByID(ctx, coll1, id, &d2); err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if d2.Version != d.Version {
			t.Errorf("versions should match")
		}
		stat, err := mdb.CollectionStats(ctx, coll1)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if stat.AvgObjSize == 0 {
			t.Errorf("Object size should be greater than zero")
		}
	})
}
