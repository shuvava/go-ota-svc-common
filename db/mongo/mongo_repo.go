package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/go-ota-svc-common/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const defaultMongoTimeout = 5 * time.Second

// BaseMongoRepository base MongoDb repository functionality
type BaseMongoRepository interface {
	db.BaseRepository
	// Database return current mongo database
	Database() *mongo.Database
	// GetCollection returns reference to mongo.Collection (table)
	GetCollection(name string) *mongo.Collection
	// InsertOne executes an insert command to insert a single document into the collection.
	InsertOne(ctx context.Context, coll *mongo.Collection, document interface{}) (string, error)
	// GetOne returns document looked up by filter, or error
	GetOne(ctx context.Context, coll *mongo.Collection, filter interface{}, document interface{}) error
	// GetOneByID returns document looked up by id, or error
	GetOneByID(ctx context.Context, coll *mongo.Collection, id string, document interface{}) error
	// Delete deletes a stored document(s) by provided filter
	Delete(ctx context.Context, coll *mongo.Collection, filter interface{}) error
	// DeleteByID deletes a stored document with provided id
	DeleteByID(ctx context.Context, coll *mongo.Collection, id string) error
	// Find returns all documents matching to the filter
	Find(ctx context.Context, coll *mongo.Collection, filter interface{}, docs interface{}) error
	// ReplaceOne replace a single document looked up by filter
	ReplaceOne(ctx context.Context, coll *mongo.Collection, filter interface{}, document interface{}) error
	// UpdateOne updates a fields in single document looked up by filter
	UpdateOne(ctx context.Context, coll *mongo.Collection, filter interface{}, update interface{}) error
	// Count returns count of documents looked up by filter
	Count(ctx context.Context, coll *mongo.Collection, filter interface{}) (int64, error)
	// CollectionStats returns general statistics about mongodb collection
	CollectionStats(ctx context.Context, coll *mongo.Collection) (*CollectionStats, error)
}

// DBResult DB result from custom queries
type DBResult map[string]interface{}

// Db service managing connection to MongoDb instance
type Db struct {
	client   *mongo.Client
	log      logger.Logger
	database string
	Timeout  time.Duration
	BaseMongoRepository
}

// CollectionStats collection statistics
type CollectionStats struct {
	// Size is collection size in bytes
	Size int64 `json:"size"`
	//AvgObjSize is Average object size in bytes
	AvgObjSize int `json:"avgObjSize"`
	//StorageSize is storage size in bytes
	StorageSize int64 `json:"storageSize"`
	//TotalIndexSize isTotal index size in bytes
	TotalIndexSize int64 `json:"totalIndexSize"`
}

// NewMongoDB create a new Db instance, with the connection URI provided
func NewMongoDB(ctx context.Context, lgr logger.Logger, connectString string) (*Db, error) {
	log := lgr.SetOperation("Db").WithContext(ctx)
	cs, err := connstring.ParseAndValidate(connectString)
	if err != nil {
		log.WithError(err).
			Error("Connection string validation failed")
		return nil, apperrors.NewAppError(
			apperrors.ErrorDbConnection,
			fmt.Sprintf("Connection string validation failed (%v)", err))
	}

	ctxConnect, cancel := context.WithTimeout(ctx, defaultMongoTimeout)
	defer cancel()
	client, err := mongo.Connect(ctxConnect, options.Client().ApplyURI(connectString))
	if err != nil {
		return nil, apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbConnection,
			"Creating NewClient failed", err)
	}

	inst := Db{
		client:   client,
		database: cs.Database,
		log:      log,
		Timeout:  defaultMongoTimeout,
	}
	return &inst, nil
}

// Disconnect close sockets to DB
func (db *Db) Disconnect(ctx context.Context) error {
	log := db.log.WithContext(ctx)
	ctxDisc, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()
	if err := db.client.Disconnect(ctxDisc); err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Disconnect from DB failed", err)
	}
	return nil
}

// Ping check connection to database
func (db *Db) Ping(ctx context.Context) error {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	ctxPing, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()
	return db.client.Ping(ctxPing, readpref.Primary())
}

// Database return current mongo database
func (db *Db) Database() *mongo.Database {
	return db.client.Database(db.database)
}

// GetCollection returns reference to mongo.Collection (table)
func (db *Db) GetCollection(name string) *mongo.Collection {
	return db.Database().Collection(name)
}

// InsertOne executes an insert command to insert a single document into the collection.
func (db *Db) InsertOne(ctx context.Context, coll *mongo.Collection, document interface{}) (string, error) {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	ctxIns, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()

	// Attempt to persist a new document
	res, err := coll.InsertOne(ctxIns, document)
	if err != nil {
		return "", apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Failed to add new DB record", err)
	}

	// Return the newly generated object ID of the persisted document
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Count returns count of documents looked up by filter
func (db *Db) Count(ctx context.Context, coll *mongo.Collection, filter interface{}) (int64, error) {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())

	ctxCnt, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()

	cnt, err := coll.CountDocuments(ctxCnt, filter)
	if err != nil {
		return 0, apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Failed to get count of DB records", err)
	}

	return cnt, nil
}

// GetOne returns document looked up by filter, or error
func (db *Db) GetOne(ctx context.Context, coll *mongo.Collection, filter interface{}, document interface{}) error {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())

	ctxGet, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()

	// Initialize a projection
	projection := bson.M{}

	opt := options.FindOne().SetProjection(projection)
	err := coll.FindOne(ctxGet, filter, opt).Decode(document)
	if err == nil {
		return nil
	}

	// If document not found, return error to indicate this
	if err == mongo.ErrNoDocuments {
		return apperrors.NewAppError(apperrors.ErrorDbNoDocumentFound, "document not found")
	}
	// Otherwise, return the provided error
	return apperrors.CreateErrorAndLogIt(log,
		apperrors.ErrorDbOperation,
		"Failed to get count of DB records", err)
}

// GetOneByID returns document looked up by id, or error
func (db *Db) GetOneByID(ctx context.Context, coll *mongo.Collection, id string, document interface{}) error {
	log := db.log.WithContext(ctx)
	oid, err := parseObjectID(id)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Invalid object ID", err)
	}

	// Compose a Filter matching the provided document ID
	filter := bson.D{primitive.E{Key: "_id", Value: oid}}

	return db.GetOne(ctx, coll, filter, document)
}

// Delete deletes a stored document(s) looked up by provided filter
func (db *Db) Delete(ctx context.Context, coll *mongo.Collection, filter interface{}) error {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())

	ctxDel, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()

	// Try to delete asset from database
	result, err := coll.DeleteMany(ctxDel, filter)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Failed no delete record from DB", err)
	}

	// If no asset was deleted, then an asset with this particular ID was not found
	if result.DeletedCount < 1 {
		return apperrors.NewAppError(
			apperrors.ErrorDbNoDocumentFound,
			"document not found")
	}

	// Otherwise, deletion was successful
	return nil
}

// DeleteByID deletes a stored document
func (db *Db) DeleteByID(ctx context.Context, coll *mongo.Collection, id string) error {
	log := db.log.WithContext(ctx)
	oid, err := parseObjectID(id)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Invalid object ID", err)
	}

	// Define a filter for this specific object ID
	filter := bson.D{primitive.E{Key: "_id", Value: oid}}

	return db.Delete(ctx, coll, filter)
}

// Find returns all documents matching to the filter
func (db *Db) Find(ctx context.Context, coll *mongo.Collection, filter interface{}, docs interface{}) error {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())

	ctxFind, cancelFind := context.WithTimeout(ctx, db.Timeout)
	defer cancelFind()

	// Initialize a projection
	projection := bson.M{}

	opt := options.Find().SetProjection(projection)
	cur, err := coll.Find(ctxFind, filter, opt)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Failed to find DB records", err)
	}

	ctxCur, cancelCur := context.WithTimeout(ctx, db.Timeout)
	defer cancelCur()
	err = cur.All(ctxCur, docs)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Failed to fetch DB records", err)
	}

	return nil
}

// ReplaceOne replace a single document looked up by filter
func (db *Db) ReplaceOne(ctx context.Context, coll *mongo.Collection, filter interface{}, document interface{}) error {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())

	ctxUpd, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()

	res, err := coll.ReplaceOne(ctxUpd, filter, document)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Failed to replace DB record", err)
	}
	if res.ModifiedCount != 1 {
		return apperrors.NewAppError(
			apperrors.ErrorDbNoDocumentFound,
			"document not found")
	}

	return nil
}

// Aggregate execute custom aggregate query
func (db *Db) Aggregate(ctx context.Context, coll *mongo.Collection, pipe interface{}, opts *options.AggregateOptions, documents interface{}) error {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())

	ctxAgg, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()

	if opts == nil {
		opts = options.Aggregate()
	}

	data, err := coll.Aggregate(ctxAgg, pipe, opts)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Failed to run DB query", err)
	}
	err = data.All(ctxAgg, documents)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"failed to decode results", err)
	}

	return nil
}

// UpdateOne updates a fields in single document looked up by filter
func (db *Db) UpdateOne(ctx context.Context, coll *mongo.Collection, filter interface{}, update interface{}) error {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())

	ctxUpd, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()

	res, err := coll.UpdateOne(ctxUpd, filter, update)
	if err != nil {
		return apperrors.CreateErrorAndLogIt(log,
			apperrors.ErrorDbOperation,
			"Failed to update DB record", err)
	}
	if res.MatchedCount == 0 {
		return apperrors.NewAppError(
			apperrors.ErrorDbNoDocumentFound,
			"document not found")
	}

	log.WithField("ModifiedCount", res.ModifiedCount).
		WithField("UpsertedCount", res.UpsertedCount).
		WithField("MatchedCount", res.MatchedCount).
		Debug("Update completed")

	return nil
}

// CollectionStats returns general statistics about mongodb collection
func (db *Db) CollectionStats(ctx context.Context, coll *mongo.Collection) (*CollectionStats, error) {
	log := db.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())

	ctxSts, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()
	result := coll.Database().RunCommand(ctxSts, bson.M{"collStats": coll.Name()})
	var doc CollectionStats
	if err := result.Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// parseObjectID is a helper to parse a string assetID into a MongoDB-format ObjectID
func parseObjectID(assetID string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return oid, nil
}
