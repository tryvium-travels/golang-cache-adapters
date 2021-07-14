// Copyright 2021 Tryvium Travels LTD
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbcacheadapters_test

import (
	ctx "context"
	"log"

	"github.com/stretchr/testify/mock"
	mongo "go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClientMock struct {
	mock.Mock
	*mongo.Client
}

func (mcm *mongoClientMock) StartSession() (mongo.Session, error) {
	args := mcm.Called()

	return args.Get(0).(mongo.Session), args.Error(1)
}

type mongoSessionMock struct {
	mock.Mock
	mongo.Session
}

func newMockMongoClientAdapter() *mongoClientMock {
	aaa := new(mongoClientMock)
	if aaa == nil {
		log.Fatalln("new mongoClientMock is wierdly nil")
	}
	clientMongo, err := mongo.Connect(ctx.Background(), options.Client().ApplyURI(localMongoDBServer.URI()))
	if err != nil {
		log.Fatalln("client mongo is not created")
	}
	aaa.Client = clientMongo
	if aaa.Client == nil {
		log.Fatalln("new Client is wierdly nil")
	}
	return aaa
}

func newMockMongoSessionAdapter() *mongoSessionMock {
	return new(mongoSessionMock)
}
