package mongodb

import (
	"os"
	"testing"
)

var dbClient MongoStore

func TestMain(m *testing.M) {
	var err error
	dbClient, err = NewMongoDBClient()
	if err != nil {
		log.Error(err)
	}
	exitVal := m.Run()
	dbClient.Close()

	os.Exit(exitVal)
}
