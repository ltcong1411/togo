package mongodb

import (
	"togo/models"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserByUserName ...
func (m *MongoClient) GetUserByUserName(username string) (user *models.User, err error) {
	filter := bson.M{
		"username": username,
	}

	err = m.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, errors.Wrapf(err, "FindOne failed - username: %v", username)
	}

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return
}

// InsertUser ...
func (m *MongoClient) InsertUser(user *models.User) (err error) {
	result, err := m.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return errors.Wrapf(err, "InsertOne failed - user: %+v", user)
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return
}
