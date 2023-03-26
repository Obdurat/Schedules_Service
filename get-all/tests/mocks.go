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

func (m *MockCollectionFail) Find(ctx context.Context, document interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return nil, errors.New("MOCK MONGO ERROR: Find failed")
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

func (m *MockCollectionSuccess) Find(ctx context.Context, document interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return mongo.NewCursorFromDocuments([]interface{}{bson.D{
		{Key: "_id", Value: "641b551ad15248336f0aad6a"},
		{Key: "client_id", Value: "abduhlaziz"},
		{Key: "service", Value: "641afadff6872fffc607baef"},
		{Key: "price", Value: 80},
		{Key: "date", Value: "2019"},
		{Key: "finished", Value: "true"},
		}}, nil, nil)
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