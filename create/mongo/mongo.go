package mongo

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	Client *mongo.Client
	Ctx context.Context
	CtxCancel context.CancelFunc
	ClientOptions *options.ClientOptions
}

type ICollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

type IRepository interface {
	Ping()
	Collection(cn string) ICollection
	Close()
}

func new() IRepository {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	options := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, options); if err != nil {
		logrus.Fatalf("%v", err)
		panic(err)
	}
	return &Repository{
		Client: client,
		Ctx: ctx,
		CtxCancel: cancel,
		ClientOptions: options,
	}
}

func (r *Repository) Ping() {
	logrus.Warn("Pinging MongoDB ...")
	if er := r.Client.Ping(r.Ctx, r.ClientOptions.ReadPreference); er != nil {
		logrus.Fatal("Mongo Database Unavailable")
		panic(er)
	}
	logrus.Infof("MongoDB connection established")
}

func (r *Repository) Collection(cn string) ICollection {
	col := r.Client.Database(cn).Collection("schedules")
	return col
}

func (r *Repository) Close() {
	r.CtxCancel()
}

var Repo = new()
