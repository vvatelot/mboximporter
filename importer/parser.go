package importer

import (
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"sync"

	"github.com/grokify/html-strip-tags-go"
	"github.com/jiphex/mbox"
)

type Parser struct {
	Filename string
	Body     bool
	Headers  bool
}

func NewParser(path string, body bool, headers bool) *Parser {
	return &Parser{
		Filename: path,
		Body:     body,
		Headers:  headers,
	}
}

func (m *Parser) ReadMessages(wg sync.WaitGroup) chan *Mail {
	queue := make(chan *Mail, 20)

	// todo debug param to flags
	messages, err := mbox.ReadFile(m.Filename, false)

	if err != nil {
		log.Fatal("unable to open file")
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		for _, message := range messages {
			tmp, err := m.parseMessage(message)

			if err != nil {
				continue
			}

			queue <- tmp
		}

		log.Println("Total: ", len(messages))

		close(queue)
	}()

	return queue
}

// https://github.com/grokify/html-strip-tags-go
func (m *Parser) parseMessage(msg *mail.Message) (*Mail, error) {

	var sender []string
	var recipients []string
	var subject []string
	var labels []string
	var message = &Mail{}

	// Headers
	headers := make([]string, len(msg.Header))
	contentType := "plain/text"

	i := 0
	for k, v := range msg.Header {
		// Specific header
		if k == "From" {
			sender = m.decodeHeaders(v)
		} else if k == "To" {
			recipients = m.decodeHeaders(v)
		} else if k == "Subject" {
			subject = m.decodeHeaders(v)
		} else if k == "Content-Type" {
			contentType = v[0]
		} else if k == "X-Gmail-Labels" {
			labels = m.splitStrings(m.decodeHeaders(v))
		}

		stringValue := k + ": " + msg.Header.Get(k)
		headers[i] = stringValue
		i++
	}

	if m.Headers == true {
		message.Headers = headers
	}

	if m.Body == true {
		// Body
		// Creates a reader.
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			return nil, err
		}

		// Reads the body
		reader := multipart.NewReader(msg.Body, params["boundary"])

		finalBody := ""
		if strings.HasPrefix(mediaType, "multipart/") {
			for {
				p, err := reader.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}
				slurp, err := ioutil.ReadAll(p)
				if err != nil {
					return nil, err
				}
				finalBody += string(slurp)
			}
		} else {
			txt, err := ioutil.ReadAll(msg.Body)
			if err != nil {
				log.Fatal(err)
			}
			finalBody += string(txt)
		}

		message.Body = strip.StripTags(finalBody)
	}

	// Reads the date
	date, err := msg.Header.Date()
	if err != nil {
		return nil, err
	}

	message.Sender = sender
	message.Recipients = recipients
	message.Date = date
	message.Subject = subject
	message.Labels = labels

	return message, nil

}

func (m *Parser) decodeHeaders(str []string) []string {
	dec := new(mime.WordDecoder)

	var decoded []string

	for _, v := range str {
		tmp, _ := dec.DecodeHeader(v)
		decoded = append(decoded, tmp)
	}
	return decoded
}

func (m *Parser) splitStrings(str []string) []string {
	var stack []string

	for _, v := range str {
		tmp := strings.Split(v, ",")

		for _, v2 := range tmp {
			stack = append(stack, v2)
		}
	}

	return stack
}
