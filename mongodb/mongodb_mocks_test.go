package mongodbcacheadapters_test

import (
	"context"

	"github.com/stretchr/testify/mock"
	mongodbcacheadapters "github.com/tryvium-travels/golang-cache-adapters/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockMongoCollection struct {
	mock.Mock
	mongodbcacheadapters.MongoCollection
	mockFind   bool
	mockInsert bool
	mockUpdate bool
	mockDelete bool
}

func (mmc *mockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	if mmc.MongoCollection != nil && !mmc.mockFind {
		return mmc.MongoCollection.FindOne(ctx, filter, opts...)
	}

	options := []interface{}{ctx, filter}

	for _, opt := range opts {
		options = append(options, opt)
	}

	args := mmc.Called(options...)

	mockMongoResult, _ := args.Get(0).(*mongo.SingleResult)
	return mockMongoResult
}

func (mmc *mockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	if mmc.MongoCollection != nil && !mmc.mockInsert {
		return mmc.MongoCollection.InsertOne(ctx, document, opts...)
	}

	options := []interface{}{ctx, document}

	for _, opt := range opts {
		options = append(options, opt)
	}

	args := mmc.Called(options...)

	mockMongoResult, _ := args.Get(0).(*mongo.InsertOneResult)
	return mockMongoResult, args.Error(1)
}

func (mmc *mockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if mmc.MongoCollection != nil && !mmc.mockUpdate {
		return mmc.MongoCollection.UpdateOne(ctx, filter, update, opts...)
	}

	options := []interface{}{ctx, filter, update}

	for _, opt := range opts {
		options = append(options, opt)
	}

	args := mmc.Called(options...)

	mockMongoResult, _ := args.Get(0).(*mongo.UpdateResult)
	return mockMongoResult, args.Error(1)
}

func (mmc *mockMongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if mmc.MongoCollection != nil && !mmc.mockDelete {
		return mmc.MongoCollection.DeleteOne(ctx, filter, opts...)
	}

	options := []interface{}{ctx, filter}

	for _, opt := range opts {
		options = append(options, opt)
	}

	args := mmc.Called(options...)

	mockMongoResult, _ := args.Get(0).(*mongo.DeleteResult)
	return mockMongoResult, args.Error(1)
}

func newMockMongoCollection(collection mongodbcacheadapters.MongoCollection, mockFind bool, mockInsert bool, mockUpdate bool, mockDelete bool) *mockMongoCollection {
	return &mockMongoCollection{
		MongoCollection: collection,
		mockFind:        mockFind,
		mockInsert:      mockInsert,
		mockUpdate:      mockUpdate,
		mockDelete:      mockDelete,
	}
}
