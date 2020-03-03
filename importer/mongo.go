package importer

import (
	"log"

	"gopkg.in/mgo.v2"
)

// Mongo connection
type Mongo struct {
	session    *mgo.Session
	collection string
	dbname     string
}

// NewConnection at MongoDB
func NewConnection(url string, dbname string, collection string) *Mongo {
	m := new(Mongo)
	session, err := mgo.Dial(url)

	if err != nil {
		panic(err)
	}

	m.session = session
	m.collection = collection
	m.dbname = dbname

	return m
}

func (m *Mongo) getCollection() *mgo.Collection {
	return m.session.DB(m.dbname).C(m.collection)
}

// BulkInsert into collection
func (m *Mongo) BulkInsert(messages []Mail) {
	bulk := m.getCollection().Bulk()
	bulk.Unordered()

	for _, message := range messages {
		bulk.Insert(message)
	}

	_, err := bulk.Run()

	if err != nil {
		log.Fatal(err)
	}
}

// Init collection (drop all data)
func (m *Mongo) Init() {
	_ = m.getCollection().DropCollection()
}

// Close connection
func (m *Mongo) Close() {
	m.session.Close()
}
