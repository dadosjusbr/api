package db

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	monthResultsCollection = "month-results"
)

//MonthResults is a data model of the results of one month parsing
type MonthResults struct {
	Month           int
	Year            int
	SpreadsheetsURL string
	DatapackageURL  string
	Success         bool
}

//Client manages all iteractions with mongodb
type Client struct {
	client *mongo.Client
	dbName string
}

//ErrDocNotFound error returned when no document is found in a query
var ErrDocNotFound = errors.New("no documents in result")

//NewClient returns an db connection instance that can be used for CRUD opetations
func NewClient(url, dbName string) (*Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(url)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return &Client{client, dbName}, nil
}

//SaveMonthResults save month results
func (db *Client) SaveMonthResults(mr MonthResults) error {
	// Get a handle for your collection
	collection := db.getMonthCollection()

	// Insert a single document
	insertResult, err := collection.InsertOne(context.TODO(), mr)
	if err != nil {
		return err
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil
}

//GetMonthResults retrieve the specified month information from the DB
func (db *Client) GetMonthResults(month, year int) (MonthResults, error) {
	filter := bson.D{
		{Key: "month", Value: month},
		{Key: "year", Value: year},
		{Key: "success", Value: true},
	}
	var result MonthResults

	err := db.getMonthCollection().FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return MonthResults{}, ErrDocNotFound
		}
		return MonthResults{}, err
	}
	return result, nil
}

//ProcessedMonth represents the information from a processed month
type ProcessedMonth struct {
	Month int
	Year  int
}

//GetProcessedMonths retrieve a list with all processed months sorted in chronological order
func (db *Client) GetProcessedMonths() ([]ProcessedMonth, error) {
	var results []ProcessedMonth
	collection := db.getMonthCollection()

	query := bson.D{{Key: "success", Value: true}}

	options := options.FindOptions{}
	options.Sort = bson.D{{Key: "year", Value: 1}, {Key: "month", Value: 1}}
	options.Projection = bson.D{{Key: "month", Value: 1}, {Key: "year", Value: 1}}

	cursor, err := collection.Find(
		context.TODO(),
		query,
		&options,
	)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var elem ProcessedMonth
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	return results, nil
}

//CloseConnection closes the opened connetion to mongodb
func (db *Client) CloseConnection() error {
	err := db.client.Disconnect(context.TODO())

	if err != nil {
		return err
	}

	fmt.Println("Connection to MongoDB closed.")
	return nil
}

func (db *Client) getMonthCollection() *mongo.Collection {
	return db.client.Database(db.dbName).Collection(monthResultsCollection)
}