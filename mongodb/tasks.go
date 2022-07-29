package mongodb

import (
	"fmt"
	"togo/models"
	"togo/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// InsertTask ...
func (m *MongoClient) InsertTask(task *models.Task) (err error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := m.Client.StartSession()
	if err != nil {
		return
	}

	defer session.EndSession(ctx)

	callback := func(sc mongo.SessionContext) (interface{}, error) {
		// get user
		user := models.User{}
		objUserID, _ := primitive.ObjectIDFromHex(task.UserID)

		err = m.UserCollection.FindOne(sc, bson.M{"_id": objUserID}).Decode(&user)
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, err
		}

		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user does not exist")
		}

		// get all today task
		filter := bson.M{
			"user_id":    task.UserID,
			"created_at": bson.M{"$gt": utils.GetStartToday()},
		}

		cur, err := m.TaskCollection.Find(sc, filter, options.Find())
		if err != nil {
			return nil, err
		}

		defer cur.Close(sc)

		tasks := []models.Task{}
		// Decode cursor
		for cur.Next(sc) {
			result := models.Task{}
			if err = cur.Decode(&result); err != nil {
				return nil, err
			}
			tasks = append(tasks, result)
		}
		if err = cur.Err(); err != nil {
			return nil, err
		}

		if len(tasks) >= user.DailyTaskLimit {
			return nil, fmt.Errorf("task limit reached")
		}

		// insert task
		result, err := m.TaskCollection.InsertOne(sc, task)
		if err != nil {
			return nil, err
		}

		task.ID = result.InsertedID.(primitive.ObjectID)

		return result, err
	}

	_, err = session.WithTransaction(ctx, callback, txnOpts)
	if err != nil {
		return
	}

	return
}
