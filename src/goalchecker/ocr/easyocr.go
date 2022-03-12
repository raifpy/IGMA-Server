package ocr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var _ OCR = &EasyOcr{} // interface Kontrol!

func init() {
	rand.Seed(time.Now().Unix())
}

type Options struct {
	Path string
}

type EasyOcr struct {
	Options Options
}

func NewEasyOcr(o Options) (ecr *EasyOcr, err error) {
	ecr = &EasyOcr{
		Options: o,
	}

	_, err = exec.LookPath(o.Path)
	return
}

func (e *EasyOcr) Ocr(ctx context.Context, path string) (ouput string, err error) {

	yanit, err := exec.CommandContext(ctx, e.Options.Path, path).CombinedOutput()
	return string(yanit), err
}

func (e *EasyOcr) Text(ctx context.Context, reader io.Reader) (string, error) {
	path := fmt.Sprintf("/tmp/ocr_media_%d", rand.Intn(9999))
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer os.Remove(path)
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}

	file.Close()

	yanit, err := e.Ocr(ctx, path)
	if err != nil {
		return "", errors.New(yanit)
	}

	return yanit, nil
}
