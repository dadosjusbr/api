package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//collection is a private interface to create a mongo's ReplaceOne method and their signatures to be used and tested.
type collection interface {
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
}

// Errors raised by package storage.
var (
	ErrNothingFound = fmt.Errorf("There is no document with this parameters")
)

//DBClient is a mongodb Client instance
type DBClient struct {
	mgoClient      *mongo.Client
	dbName         string
	monthlyInfoCol string
	agencyCol      string
	col            collection
}

//NewDBClient instantiates a mongo new client, but will not connect to the specified URL. Please use Client.Connect before using the client.
func NewDBClient(url, dbName, monthlyInfoCol, agencyCol string) (*DBClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	return &DBClient{mgoClient: client, dbName: dbName, monthlyInfoCol: monthlyInfoCol, agencyCol: agencyCol}, nil
}

//Connect establishes a connection to MongoDB using the previously specified URL
func (c *DBClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.mgoClient.Connect(ctx); err != nil {
		return fmt.Errorf("error connection with mongo:%q", err)
	}
	return nil
}

//Disconnect closes the connections to MongoDB. It does nothing if the connection had already been closed.
func (c *DBClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	c.col = nil
	return c.mgoClient.Disconnect(ctx)
}

// GetDataForFirstScreen GetDataForFirstScreen
func (c *DBClient) GetDataForFirstScreen(uf string, year int) ([]Agency, map[string][]AgencyMonthlyInfo, error) {
	allAgencies, err := c.GetAgencies(uf)
	if err != nil {
		return nil, nil, fmt.Errorf("GetDataForFirstScreen() error: %q", err)
	}
	result, err := c.GetMonthlyInfo(allAgencies, year)
	if err != nil {
		return nil, nil, fmt.Errorf("GetDataForFirstScreen() error: %q", err)
	}
	return allAgencies, result, nil
}

//GetAgencies Return UF Agencies
func (c *DBClient) GetAgencies(uf string) ([]Agency, error) {
	c.Collection(c.agencyCol)
	resultAgencies, err := c.col.Find(context.TODO(), bson.D{{Key: "uf", Value: uf}}, nil)
	if err != nil {
		return nil, fmt.Errorf("Find error in getAgencies %v", err)
	}
	var allAgencies []Agency
	resultAgencies.All(context.TODO(), &allAgencies)
	if err := resultAgencies.Err(); err != nil {
		return nil, fmt.Errorf("Error in result %v", err)
	}
	return allAgencies, nil
}

//GetMonthlyInfo return summarized monthlyInfo for each agency in agencies in a specific year
func (c *DBClient) GetMonthlyInfo(agencies []Agency, year int) (map[string][]AgencyMonthlyInfo, error) {
	var result = make(map[string][]AgencyMonthlyInfo)
	c.Collection(c.monthlyInfoCol)
	findOptions := options.Find()
	for _, agency := range agencies {
		resultMonthly, err := c.col.Find(context.TODO(), bson.D{{Key: "aid", Value: agency.ID}, {Key: "year", Value: year}},
			findOptions.SetProjection(bson.D{{Key: "aid", Value: ""}, {Key: "year", Value: ""}, {Key: "month", Value: ""}, {Key: "summary", Value: ""}}))
		if err != nil {
			return nil, fmt.Errorf("Error in GetMonthlyInfo %v", err)
		}
		var mr []AgencyMonthlyInfo
		resultMonthly.All(context.TODO(), &mr)
		result[agency.ID] = mr
	}

	return result, nil
}

//GetDataForSecondScreen Search if DB has a match for filters
func (c *DBClient) GetDataForSecondScreen(month int, year int, agency string) (*AgencyMonthlyInfo, error) {
	c.Collection(c.monthlyInfoCol)
	var resultMonthly AgencyMonthlyInfo
	err := c.col.FindOne(context.TODO(), bson.D{{Key: "aid", Value: agency}, {Key: "year", Value: year}, {Key: "month", Value: month}}).Decode(&resultMonthly)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		return nil, ErrNothingFound
	}
	return &resultMonthly, nil
}

//Collection Changes active collection
func (c *DBClient) Collection(collectionName string) {
	c.col = c.mgoClient.Database(c.dbName).Collection(collectionName)
}
