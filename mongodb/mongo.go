package mongodb

import (
	"log"

	"gopkg.in/mgo.v2"

	"github.com/rpsl/mboximporter/importer"
)

type Mongo struct {
	session    *mgo.Session
	database   *mgo.Database
	collection string
	dbname     string
}

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

func (m *Mongo) BulInsert(messages []importer.Mail) {
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

func (m *Mongo) Init() {
	err := m.getCollection().DropCollection()

	if err != nil {
		// log.Fatal(err)
	}
}

func (m *Mongo) Close() {
	m.session.Close()
}
