package main

import (
	"log"
	"os"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/gen2brain/beeep"
)

var (
	allIds   []uint32
	logger   *log.Logger
	email    string
	password string
)

type EmailHeader struct {
	Sender  string
	Subject string
}

func setupLogger() *os.File {
	f, err := os.OpenFile("log/outlook.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Println(err)
	}

	logger = log.New(f, "", log.LstdFlags)
	return f
}

func addToDatabase(emailMap map[uint32]EmailHeader) {
	for id, _ := range emailMap {
		allIds = append(allIds, id)
	}
}

func clearDatabase() {
	allIds = nil
}

func contains(list []uint32, element uint32) bool {
	for _, v := range list {
		if element == v {
			return true
		}
	}
	return false
}

func sendNotifications(emailMap map[uint32]EmailHeader) {
	for id, header := range emailMap {
		if contains(allIds, id) {
			continue
		} else {
			err := beeep.Notify(header.Sender, header.Subject, "assets/outlook.ico")
			if err != nil {
				logger.Fatal(err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func getNameAndEmail(msg *imap.Message) string {
	from := msg.Envelope.From[0]
	return from.PersonalName + " <" + from.MailboxName + "@" + from.HostName + ">"
}

func connectToMailServer() *client.Client {
	logger.Println("Connecting to server...")

	// Connect to server
	server := "outlook.office365.com:993"
	c, err := client.DialTLS(server, nil)
	if err != nil {
		logger.Fatal(err)
	}
	return c
}

func getUnreadEmails(c *client.Client) map[uint32]EmailHeader {
	// Select INBOX
	readonly := true
	_, err := c.Select("INBOX", readonly)
	if err != nil {
		logger.Fatal(err)
	}

	// Set search criteria
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}

	ids, err := c.Search(criteria)
	if err != nil {
		logger.Fatal(err)
	}

	var headers []EmailHeader
	emailMap := make(map[uint32]EmailHeader)
	if len(ids) > 0 {
		logger.Println("# IDs found:", len(ids), "--", ids)
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()

		for msg := range messages {
			h := EmailHeader{
				Sender:  getNameAndEmail(msg),
				Subject: msg.Envelope.Subject,
			}
			headers = append(headers, h)
		}

		if err := <-done; err != nil {
			logger.Fatal(err)
		}

	}

	if len(ids) != len(headers) {
		logger.Fatal("Length of Email IDs and Email Headers are mismatched.")
	}
	for i := 0; i < len(ids); i++ {
		emailMap[ids[i]] = headers[i]
	}
	return emailMap
}

func checkEmail() {
	f := setupLogger()
	defer f.Close()

	// Connect to Outlook server
	c := connectToMailServer()

	// Login
	err := c.Login(email, password)
	if err != nil {
		logger.Fatal(err)
	}
	defer c.Logout()

	// Get Unseen Emails
	emailMap := getUnreadEmails(c)

	if len(emailMap) > 0 {
		sendNotifications(emailMap)
		addToDatabase(emailMap)
	}
}

func main() {
	// Set these with `export OUTLOOK_EMAIL=...`
	email = os.Getenv("OUTLOOK_EMAIL")
	password = os.Getenv("OUTLOOK_PASSWORD")

	interval := 30 * time.Second
	resetTime := 24 * time.Hour
	for i := 0; ; i++ {
		checkEmail()
		time.Sleep(interval)
		// Empty database every 24 hours
		if time.Duration(i)*interval > resetTime {
			clearDatabase()
			i = 0
		}
	}
}
