package mongodb

import (
	"context"

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
	DB             *mongo.Database
	UserCollection *mongo.Collection
	TaskCollection *mongo.Collection
}

// MongoStore interface
type MongoStore interface {
	// users
	InsertUser(user *models.User) (err error)
	GetUserByUserName(username string) (user *models.User, err error)

	// tasks
	InsertTask(task *models.Task) (err error)

	Close()
}

// NewMongoDBClient ...
func NewMongoDBClient() (MongoStore, error) {
	mongoAddress := "mongodb://" + config.Values.Mongo.Host + ":" + config.Values.Mongo.Port

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
		DB:             mongoCli.Database(config.Values.Mongo.DB),
		UserCollection: mongoCli.Database(config.Values.Mongo.DB).Collection(config.Values.Mongo.Collection.User),
		TaskCollection: mongoCli.Database(config.Values.Mongo.DB).Collection(config.Values.Mongo.Collection.Task),
	}

	return mongo, nil
}

// Close connection
func (m *MongoClient) Close() {
	if m.Client != nil {
		m.Client.Disconnect(ctx)
	}
}
