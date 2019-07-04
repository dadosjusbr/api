package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// INSTRUCTIONS:
// * To run this test you first need to start MongoDB
// $ mkdir /tmp/rmtest
// $ mongod --port 27017 --dbpath=/tmp/rmtest --httpinterface --rest &

const dbName = "testDB"
const urlString = "mongodb://localhost:27017/"

var (
	dbClient    *Client       // client which we must test against
	queryClient *mongo.Client // client which we should use to validate tests
)

func TestMain(m *testing.M) {
	var err error
	dbClient, err = NewClient(urlString, dbName)
	if err != nil {
		log.Fatalf("error creating dbClient:%q", err)
	}
	queryClient, err = mongo.NewClient(options.Client().ApplyURI(urlString))
	if err != nil {
		log.Fatalf("error creating queryClient:%q", err)
	}
	if err := queryClient.Connect(context.Background()); err != nil {
		log.Fatalf("error connecting queryClient:%q", err)
	}
	retCode := m.Run()
	dbClient.CloseConnection()
	queryClient.Database(dbName).Drop(context.Background())
	queryClient.Disconnect(context.Background())
	os.Exit(retCode)
}

func TestClient_SaveMonthResults(t *testing.T) {
	t.Run("FirstInsert", func(t *testing.T) {
		is := is.New(t)
		mr := MonthResults{Year: 2010, Month: 10}
		is.NoErr(dbClient.SaveMonthResults(mr))

		col := queryClient.Database(dbName).Collection(monthResultsCollection)
		filter := bson.D{{Key: "_id", Value: getID(mr)}}
		defer col.DeleteOne(context.Background(), filter)
		res := col.FindOne(context.Background(), filter)
		is.NoErr(res.Err())

		var resMR MonthResults
		is.NoErr(res.Decode(&resMR))
		is.Equal(mr.Year, resMR.Year)
		is.Equal(mr.Month, resMR.Month)
	})

	t.Run("Upsert", func(t *testing.T) {
		is := is.New(t)
		mr := MonthResults{Year: 2010, Month: 10}
		is.NoErr(dbClient.SaveMonthResults(mr))
		is.NoErr(dbClient.SaveMonthResults(mr))

		col := queryClient.Database(dbName).Collection(monthResultsCollection)
		n, err := col.CountDocuments(context.Background(), bson.D{{Key: "month", Value: 10}})
		is.NoErr(err)
		is.Equal(n, int64(1))

		filter := bson.D{{Key: "_id", Value: getID(mr)}}
		defer col.DeleteOne(context.Background(), filter)
		res := col.FindOne(context.Background(), filter)
		is.NoErr(res.Err())

		var resMR MonthResults
		is.NoErr(res.Decode(&resMR))
		is.Equal(mr.Year, resMR.Year)
		is.Equal(mr.Month, resMR.Month)
	})
}
