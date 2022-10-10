package importer

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo connection
type Mongo struct {
	client     *mongo.Client
	collection string
	dbname     string
}

// NewConnection at MongoDB
func NewConnection(url string, dbname string, collection string) *Mongo {
	m := new(Mongo)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+url))

	if err != nil {
		panic(err)
	}

	m.client = client
	m.collection = collection
	m.dbname = dbname

	return m
}

func (m *Mongo) getCollection() *mongo.Collection {
	return m.client.Database(m.dbname).Collection(m.collection)
}

// BulkInsert into collection
func (m *Mongo) BulkInsert(messages []Mail) {
	bulk := m.getCollection()
	var ui []interface{}
	for _, t := range messages {
		ui = append(ui, t)
	}

	_, err := bulk.InsertMany(context.TODO(), ui)
	if err != nil {
		log.Fatal(err)
	}
}

// Init collection (drop all data)
func (m *Mongo) Init() {
	_ = m.getCollection().Drop(context.TODO())
}

// Close connection
func (m *Mongo) Close() {
	m.client.Disconnect(context.TODO())
}
