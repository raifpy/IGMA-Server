package ocr

import "io"

type Ocr interface {
	Text(reader io.Reader) (string, error)
}
