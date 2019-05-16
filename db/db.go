package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	monthResultsCollection = "month-results"
	dbName                 = "dadosjusbr"
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
}

//NewClient returns an db connection instance that can be used for CRUD opetations
func NewClient(url string) (*Client, error) {
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
	return &Client{client}, nil
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
	return db.client.Database(dbName).Collection(monthResultsCollection)
}
