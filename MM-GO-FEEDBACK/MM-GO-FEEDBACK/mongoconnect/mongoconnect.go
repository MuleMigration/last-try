package MongoConnect

import (
	"context"
	MMErr "feedback/mmerror"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBInterface interface {
	GetAppClient() (*mongo.Database, *MMErr.AppError)
	GetUserClient(database string) (*mongo.Database, *MMErr.AppError)
}

func NewMongo() MongoDB {
	return MongoDB{}
}

type MongoDB struct{}

func (db *MongoDB) GetAppClient() (*mongo.Database, *MMErr.AppError) {
	secrets := getSecrets()
	client, err := db.connectToMongo(secrets.MongoDBURI, secrets.AppDatabase)
	if err != nil {
		return nil, err
	}
	return client.Database(secrets.AppDatabase), nil
}

func (db *MongoDB) GetUserClient(database string) (*mongo.Database, *MMErr.AppError) {
	secrets := getSecrets()
	client, err := db.connectToMongo(secrets.MongoDBURI, database)
	if err != nil {
		return nil, err
	}
	return client.Database(database), nil
}

func (db *MongoDB) connectToMongo(connectionString, database string) (*mongo.Client, *MMErr.AppError) {
	clientOptions := options.Client().ApplyURI(connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, MMErr.NewUnexpectedError(err.Error())
	}
	return client, nil
}

func getSecrets() *Secrets {
	// Implement this function to get the MongoDB connection URI and Application Database name.
	// This is just a placeholder and needs to be replaced with actual implementation.
	return &Secrets{
		MongoDBURI:  "mongodb+srv://devreader:MuleMigration%40123!!@mulemigration-dev.vnsm7lg.mongodb.net/",
		AppDatabase: "ApplicationDB",
	}
}

type Secrets struct {
	MongoDBURI  string
	AppDatabase string
}
