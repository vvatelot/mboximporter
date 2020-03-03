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

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jiphex/mbox"
)

const channelSize = 20

// Parser ...
type Parser struct {
	Filename string
	Body     bool
	Headers  bool
}

// NewParser ...
func NewParser(path string, body bool, headers bool) *Parser {
	return &Parser{
		Filename: path,
		Body:     body,
		Headers:  headers,
	}
}

// ReadMessages from file return them by channel
// todo: split return type to separate methods
func (m *Parser) ReadMessages(wg *sync.WaitGroup) chan *Mail {
	queue := make(chan *Mail, channelSize)

	// todo debug param to flags
	messages, err := mbox.ReadFile(m.Filename, false)

	if err != nil {
		log.Fatal("unable to open file")
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for _, message := range messages {
			msg, err := m.parseMessage(message)

			if err != nil {
				continue
			}

			queue <- msg
		}

		close(queue)

		log.Println("Total: ", len(messages))
	}()

	return queue
}

func (m *Parser) parseMessage(msg *mail.Message) (*Mail, error) {
	message := m.parseHeader(msg)

	if m.Body {
		body, err := m.parseBody(msg)

		if err != nil {
			return nil, err
		}

		message.Body = body
	}

	// Reads the date
	date, err := msg.Header.Date()
	if err != nil {
		return nil, err
	}

	message.Date = date

	return message, nil
}

func (m *Parser) parseHeader(msg *mail.Message) *Mail {
	var (
		sender, recipients, subject, labels, headers []string

		message = &Mail{}
	)

	for k, v := range msg.Header {
		// Specific header
		switch k {
		case "From":
			sender = m.decodeHeaders(v)
		case "To":
			recipients = m.decodeHeaders(v)
		case "Subject":
			subject = m.decodeHeaders(v)
		case "X-Gmail-Labels":
			labels = m.splitStrings(m.decodeHeaders(v))
		default:
			if m.Headers {
				stringValue := k + ": " + msg.Header.Get(k)
				headers = append(headers, stringValue)
			}
		}
	}

	if len(headers) > 0 {
		message.Headers = headers
	}

	message.Sender = sender
	message.Recipients = recipients
	message.Subject = subject
	message.Labels = labels

	return message
}

func (m *Parser) parseBody(msg *mail.Message) (string, error) {
	contentType := "plain/text"

	check := msg.Header.Get("Content-Type")

	if check != "" {
		contentType = check
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
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
				return "", err
			}

			slurp, err := ioutil.ReadAll(p)

			if err != nil {
				return "", err
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

	finalBody = strip.StripTags(finalBody)

	return finalBody, nil
}

func (m *Parser) decodeHeaders(str []string) []string {
	// var stack []string
	stack := make([]string, 0, len(str))
	coder := new(mime.WordDecoder)

	for _, v := range str {
		text, err := coder.DecodeHeader(v)

		if err != nil || text == "" {
			continue
		}

		stack = append(stack, text)
	}

	return stack
}

func (m *Parser) splitStrings(str []string) []string {
	var stack []string

	for _, v := range str {
		text := strings.Split(v, ",")

		stack = append(stack, text...)
	}

	return stack
}
