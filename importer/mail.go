package importer

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Mail message struct inserted into MongoDB
type Mail struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	Headers    []string
	Sender     []string
	Recipients []string
	Labels     []string
	Subject    []string
	Date       time.Time
	Body       string
}
