package mail

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	fp "path/filepath"

	"github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/logging"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	gmailSenderEmailEnvKey = "gmail_sender_email"
)

var (
	gmailSenderEmail = os.Getenv(gmailSenderEmailEnvKey)
)

var srv *gmail.Service

// Send an email with attachement using GMail API
// https://stackoverflow.com/a/62214410
func Send(email, filepath string) error {
	if srv == nil {

		ctx := context.Background()
		b, err := os.ReadFile("credentials.json")
		if err != nil {
			log.Printf("Unable to read client secret file: %v", err)
			return err
		}

		// If modifying these scopes, delete your previously saved token.json.
		config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, gmail.GmailSendScope)
		if err != nil {
			log.Printf("Unable to parse client secret file to config: %v", err)
			return err
		}
		client, err := getClient(config)

		if err != nil {
			return err
		}

		srv, err = gmail.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			log.Printf("Unable to retrieve Gmail client: %v", err)
			return err
		}
	}

	user := "me"

	var message gmail.Message

	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	fileMIMEType := http.DetectContentType(fileBytes)
	filename := fp.Base(filepath)

	fileData := base64.StdEncoding.EncodeToString(fileBytes)

	boundary := randStr(32, "alphanum")

	messageBody := []byte("Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
		"MIME-Version: 1.0\n" +
		"to: " + email + "\n" +
		"from: " + gmailSenderEmail + "\n" +
		"subject: " + "Send To PocketBook" + "\n\n" +

		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: 7bit\n\n" +
		"" + "\n\n" +
		"--" + boundary + "\n" +

		"Content-Type: " + fileMIMEType + "; name=" + string('"') + filename + string('"') + " \n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: base64\n" +
		"Content-Disposition: attachment; filename=" + string('"') + filename + string('"') + " \n\n" +
		chunkSplit(fileData, 76, "\n") +
		"--" + boundary + "--")

	message.Raw = base64.URLEncoding.EncodeToString(messageBody)

	_, err = srv.Users.Messages.Send(user, &message).Do()
	if err != nil {
		return err
	}

	logging.Aspirador.Info(fmt.Sprintf("Email sent to %s successfuly", email))

	return nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Printf("token.json not found")
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func chunkSplit(body string, limit int, end string) string {

	var charSlice []rune

	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result string = ""

	for len(charSlice) >= 1 {

		result = result + string(charSlice[:limit]) + end

		charSlice = charSlice[limit:]

		if len(charSlice) < limit {
			limit = len(charSlice)
		}

	}

	return result

}
