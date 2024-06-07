package http

import (
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"

	"github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/logging"
	sender "github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/mail"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {

	e.POST(sendRoute, SendToPocketBook)
	e.GET(internalServerErrorRoute, ReturnInternalServerError)

	echo.NotFoundHandler = useNotFoundHandler()
}

func SendToPocketBook(ectx echo.Context) error {
	logging.Aspirador.Trace("handling send-to-pocket-book")

	request := SendToPocketBookRequest{}
	decerr := ectx.Bind(&request)

	if decerr != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while reading request body: %s", decerr))
		return InternalServerError(ectx)
	}

	if _, err := mail.ParseAddress(request.Email); err != nil {
		fmt.Printf("err: %v\n", err)

		return BadRequest(ectx)
	}

	res, err := http.Head(request.Url)

	if err != nil {
		fmt.Printf("err: %v\n", err)

		return BadRequest(ectx)
	}

	contentType := res.Header.Get("content-type")
	fileExt, err := contentTypeToFileExtension(contentType)

	if err != nil {
		fmt.Printf("err: %v\n", err)

		return BadRequest(ectx)
	}

	tmpfile, err := os.CreateTemp("", "")

	if err != nil {
		fmt.Printf("err: %v\n", err)

		return BadRequest(ectx)
	}

	filepath := fmt.Sprintf("%s-%s%s", tmpfile.Name(), request.Title, fileExt)

	err = downloadFile(filepath, request.Url)

	if err != nil {
		fmt.Printf("err: %v\n", err)

		return BadRequest(ectx)
	}

	s := sender.NewSender()
	m := sender.NewMessage("Send to PocketBook", "")
	m.To = []string{request.Email}
	m.AttachFile(filepath)

	err = s.Send(m)

	if err != nil {
		fmt.Printf("err: %v\n", err)

		return BadRequest(ectx)
	}

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
