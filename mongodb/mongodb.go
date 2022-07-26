package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"togo/config"
	"togo/logger"
	"togo/models"
)

var (
	log = logger.GetLogger("handlers")
	ctx = context.TODO()
)

// MongoClient manage all mongodb action
type MongoClient struct {
	Client         *mongo.Client
	UserCollection *mongo.Collection
}

// MongoStore interface
type MongoStore interface {
	// users
	GetUserByUserName(username string) (user *models.User, err error)
	InsertUser(user *models.User) (err error)

	Close()
}

// NewMongoDBClient ...
func NewMongoDBClient() (MongoStore, error) {
	// mongoAddress := "mongodb://" + config.Values.Mongo.Host + ":" + config.Values.Mongo.Port
	mongoAddress := fmt.Sprintf(config.Values.Mongo.Address, config.Values.Mongo.Username, config.Values.Mongo.Password, config.Values.Mongo.DB)
	log.Debug("mongoAddress: ", mongoAddress)

	mongoCli, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoAddress))
	if err != nil {
		log.Errorf("Error when try to connect to Mongodb server: %v\n", err.Error())
		return nil, err
	}

	if err := mongoCli.Ping(ctx, nil); err != nil {
		log.Errorf("Can not ping to Mongodb server: %v\n", err)
	}
	log.Info("Connected to Mongodb Server")

	mongo := &MongoClient{
		Client:         mongoCli,
		UserCollection: mongoCli.Database(config.Values.Mongo.DB).Collection(config.Values.Mongo.Collection.User),
	}

	return mongo, nil
}

// Close connection
func (m *MongoClient) Close() {
	if m.Client != nil {
		m.Client.Disconnect(ctx)
	}
}
