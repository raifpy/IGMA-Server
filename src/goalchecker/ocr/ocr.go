package ocr

import (
	"context"
	"io"
)

type OCR interface {
	Text(context.Context, io.Reader) (string, error)
}
