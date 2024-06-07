package http

type SendToPocketBookRequest struct {
	Email string `json:"email"`
	Title string `json:"title"`
	Url   string `json:"url"`
}
