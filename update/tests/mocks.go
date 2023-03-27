package tests

import (
	"context"
	"errors"

	repository "github.com/Obdurat/Schedules/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MOCK FAILURE ON MONGO DB -----------------------------------------------------------------

type MockCollectionFail struct {}

func (m *MockCollectionFail) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return mongo.NewSingleResultFromDocument(bson.D{{Key: "doc", Value: "aqui"}}, errors.New("MOCK ERROR: FindOneAndUpdate"), nil)
}

type MockRepoFail struct {}

func (m *MockRepoFail) Ping() {
	return;
}

func (m *MockRepoFail) Database(cn string) repository.ICollection {
	return &MockCollectionFail{}
}

func (m *MockRepoFail) Disconnect() {
	return
}

// MOCK SUCCESS ON MONGO DB -----------------------------------------------------------------

type MockCollectionSuccess struct {}

func (m *MockCollectionSuccess) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return mongo.NewSingleResultFromDocument(bson.D{{Key: "doc", Value: "aqui"}}, nil, nil)
}

type MockRepoSuccess struct {}

func (m *MockRepoSuccess) Ping() {
	return;
}

func (m *MockRepoSuccess) Database(cn string) repository.ICollection {
	return &MockCollectionSuccess{}
}

func (m *MockRepoSuccess) Disconnect() {
	return
}

// INSTANCE BASED ON MODE --------------------------------------------------------------------

func new(mode string) repository.IRepository {
	if mode == "MONGO_SUCCESS" {
		return &MockRepoSuccess{}
	}
	return &MockRepoFail{}
}