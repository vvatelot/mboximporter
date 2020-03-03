package cmd

import (
	"log"
	"sync"

	"github.com/rpsl/mboximporter/importer"
)

// Importer command
type Importer struct {
	Mongo    string `help:"The Mongo URI to connect to MongoDB"`
	Database string `help:"The Database name to use in MongoDB"`
	Filename string `help:"Name of the filename to import"`
	Body     bool   `help:"Parse and insert body of the emails"`
	Headers  bool   `help:"Parse and insert all headers of the emails"`
	Init     bool   `help:"Drop if exist collection and create fresh"`
}

const batchSize = 500

// NewImport command
// todo: sync names with struct
func NewImport() *Importer {
	return &Importer{
		Mongo:    "root:example@127.0.0.1",
		Database: "mbox-importer",
		Body:     false,
		Headers:  false,
		Init:     false,
	}
}

// Run ...
func (m *Importer) Run() error {
	mongo := importer.NewConnection(m.Mongo, m.Database, "mails")
	defer mongo.Close()

	if m.Init {
		mongo.Init()
	}

	var wg sync.WaitGroup

	parser := importer.NewParser(m.Filename, m.Body, m.Headers)
	queue := parser.ReadMessages(&wg)

	wg.Add(1)

	go func() {
		defer wg.Done()

		var stack []importer.Mail

		count := 0

		for message := range queue {
			count++

			stack = append(stack, *message)

			if len(stack) == batchSize {
				mongo.BulkInsert(stack)

				stack = nil

				log.Println("Worked: ", count)
			}
		}

		// Last insert
		mongo.BulkInsert(stack)
	}()

	wg.Wait()

	return nil
}
