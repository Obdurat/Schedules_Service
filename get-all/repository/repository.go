package repository

import (
	"context"
	"os"
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
	Find(ctx context.Context, document interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
}

type IRepository interface {
	Ping()
	Database(cn string) ICollection
	Disconnect()
}

func new() *Repository {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	options := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, options); if err != nil {
		logrus.Fatal(err)
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

func (r *Repository) Database(cn string) ICollection {
	col := r.Client.Database(cn).Collection("schedules")
	return col
}

func (r *Repository) Disconnect() {
	r.CtxCancel()
}

var Instance = new()
