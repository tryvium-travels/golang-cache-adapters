package mongodbcacheadapters_test

import (
	"github.com/stretchr/testify/mock"
	mongodbcacheadapters "github.com/tryvium-travels/golang-cache-adapters/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockMongoClient struct {
	mock.Mock
	mongodbcacheadapters.MongoClient
}

func (mmc *mockMongoClient) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	options := make([]interface{}, 0, len(opts))

	for _, opt := range opts {
		options = append(options, opt)
	}

	args := mmc.Called(options...)
	sessionArg, _ := args.Get(0).(mongo.Session)
	return sessionArg, args.Error(1)
}

func newMockMongoClient() *mockMongoClient {
	return new(mockMongoClient)
}
