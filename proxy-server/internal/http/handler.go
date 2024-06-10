package http

import (
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"

	"github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/logging"
	sender "github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/mail"
	"github.com/labstack/echo/v4"
)

var (
	asciiCharactersAndNumbersOnlyRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
)

func RegisterHandlers(e *echo.Echo) {

	e.POST(sendRoute, SendToPocketBook)
	e.GET(internalServerErrorRoute, ReturnInternalServerError)

	echo.NotFoundHandler = useNotFoundHandler()
}

func SendToPocketBook(ectx echo.Context) error {
	logging.Aspirador.Trace("handling send-to-pocket-book")

	breq := SendToPocketBookRequest{}
	decerr := ectx.Bind(&breq)

	if decerr != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while reading request body: %s", decerr))
		return InternalServerError(ectx)
	}

	if _, err := mail.ParseAddress(breq.Email); err != nil {
		fmt.Printf("err: %v\n", err)

		return BadRequest(ectx, InvalidEmailResponse())
	}

	if !strings.HasSuffix(breq.Email, "@pbsync.com") {
		return BadRequest(ectx, UnsupportedEmailDomainResponse())
	}

	res, err := http.Head(breq.Url)

	if err != nil {
		fmt.Printf("err: %v\n", err)

		return BadRequest(ectx, NotConnectableUrlResponse())
	}

	contentType := res.Header.Get("content-type")
	fileExt, err := contentTypeToFileExtension(strings.Split(contentType, ";")[0])

	if err != nil {
		fmt.Printf("err: %v\n", err)

		return BadRequest(ectx, UnsupportedDocumentResponse())
	}

	tmpfile, err := os.CreateTemp("", "")

	if err != nil {
		fmt.Printf("err: %v\n", err)

		return InternalServerError(ectx)
	}

	go func() {
		title := sanitizeFileName(breq.Title)
		filepath := fmt.Sprintf("%s-%s%s", tmpfile.Name(), title, fileExt)

		err = downloadFile(filepath, breq.Url)

		if err != nil {
			fmt.Printf("err: %v\n", err)

			return
		}

		err = sender.Send(breq.Email, filepath)

		if err != nil {
			fmt.Printf("err: %v\n", err)

			return
		}

		os.Remove(filepath)
	}()

	return Accepted(ectx)
}

func ReturnInternalServerError(ectx echo.Context) error {
	logging.Aspirador.Trace("Returning Internal Server Error")
	return InternalServerError(ectx)
}

func useNotFoundHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	}
}

func contentTypeToFileExtension(contentType string) (string, error) {
	var ext string
	var err error

	switch contentType {
	case "application/pdf":
		ext = ".pdf"
	case "text/html":
		ext = ".html"
	case "text/plain":
		ext = ".txt"
	default:
		err = fmt.Errorf("do not recognize %s as file extension", contentType)
	}

	return ext, err
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func sanitizeFileName(filename string) string {
	sfn := filename

	if len(sfn) > 20 {
		sfn = string(sfn[:20])
	}

	sfn = asciiCharactersAndNumbersOnlyRegex.ReplaceAllString(sfn, "")

	if len(sfn) == 0 {
		sfn = "unknown"
	}

	return sfn

}
