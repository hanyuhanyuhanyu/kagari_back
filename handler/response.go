package handler

type ContentType string

const (
	JSON = ContentType("json")
)

type Response struct {
	Status      int
	Body        any
	ContentType ContentType
}
