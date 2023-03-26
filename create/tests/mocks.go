package tests

import (
	"context"
	"errors"

	repo "github.com/Obdurat/Schedules/mongo"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MOCK FAILURE ON MONGO DB -----------------------------------------------------------------

type MockCollectionFail struct {
	mock.Mock
}

func (m *MockCollectionFail) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, errors.New("inducedError")
}

type MockRepoFail struct {
	mock.Mock
}

func (m *MockRepoFail) Ping() {
	return;
}

func (m *MockRepoFail) Collection(cn string) repo.ICollection {
	return &MockCollectionFail{}
}

func (m *MockRepoFail) Close() {
	return
}

// MOCK SUCCESS ON MONGO DB -----------------------------------------------------------------

type MockCollectionSuccess struct {
	mock.Mock
}

func (m *MockCollectionSuccess) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}

type MockRepoSuccess struct {
	mock.Mock
}

func (m *MockRepoSuccess) Ping() {
	return;
}

func (m *MockRepoSuccess) Collection(cn string) repo.ICollection {
	return &MockCollectionSuccess{}
}

func (m *MockRepoSuccess) Close() {
	return
}

func new(mode string) repo.IRepository {
	if mode == "MONGO_SUCCESS" {
		return &MockRepoSuccess{}
	}
	return &MockRepoFail{}
}