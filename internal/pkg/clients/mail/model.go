package mail

import (
	"io"
)

const (
	FieldFrom       = "From"
	FieldTo         = "To"
	FieldSubject    = "Subject"
	ContentTypeHTML = "text/html"
)

type Attachment struct {
	Filename string
	Content  io.Reader
}
