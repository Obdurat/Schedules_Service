package tests

import (
	"context"
	"errors"

	repository "github.com/Obdurat/Schedules/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MOCK FAILURE ON MONGO DB -----------------------------------------------------------------

type MockCollectionFail struct {}

func (m *MockCollectionFail) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, errors.New("MOCK ERROR: InsertOne")
}

type MockRepoFail struct {}

func (m *MockRepoFail) Ping() {
	return;
}

func (m *MockRepoFail) Collection(cn string) repository.ICollection {
	return &MockCollectionFail{}
}

func (m *MockRepoFail) Disconnect() {
	return
}

// MOCK SUCCESS ON MONGO DB -----------------------------------------------------------------

type MockCollectionSuccess struct {}

func (m *MockCollectionSuccess) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}

type MockRepoSuccess struct {}

func (m *MockRepoSuccess) Ping() {
	return;
}

func (m *MockRepoSuccess) Collection(cn string) repository.ICollection {
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