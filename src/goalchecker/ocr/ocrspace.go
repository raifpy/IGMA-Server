package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type OcrSpaceOcr struct {
	ServerURL string
	Token     string
	//Body      url.Values
	Body map[string]string

	Client *http.Client
}

// {"apikey": "5a64d478-9c89-43d8-88e3-c65de9999580"}
type Config struct {
	OCREngine string // 1 2
}

func NewOcrSpaceOcr(c ...Config) OcrSpaceOcr {
	var engine string = "2"
	if len(c) != 0 && c[0].OCREngine != "" {
		engine = c[0].OCREngine
	}

	return OcrSpaceOcr{
		ServerURL: "https://api8.ocr.space/parse/image",
		//Client: &http.Client{},
		Token:  "5a64d478-9c89-43d8-88e3-c65de9999580",
		Client: http.DefaultClient,
		/*Body: url.Values{
			"url":                          {""},
			"language":                     {"tur"},
			"isOverlayRequired":            {"true"},
			"FileType":                     {".Auto"},
			"IsCreateSearchablePDF":        {"false"},
			"isSearchablePdfHideTextLayer": {"true"},
			"detectOrientation":            {"false"},
			"isTable":                      {"false"},
			"scale":                        {"false"},
			"OCREngine":                    {"1"},
			"detectCheckbox":               {"false"},
			"checkboxTemplate":             {"0"},
		},*/
		Body: map[string]string{
			"url": "",
			//"language": "tur",
			"language":                     "eng",
			"isOverlayRequired":            "true",
			"FileType":                     ".Auto",
			"IsCreateSearchablePDF":        "false",
			"isSearchablePdfHideTextLayer": "true",
			"detectOrientation":            "false",
			"isTable":                      "false",
			"scale":                        "false",
			//"OCREngine":                    "1",
			"OCREngine": engine,
			//"OCREngine":        "2",
			"detectCheckbox":   "false",
			"checkboxTemplate": "0",
		},
	}
}

type OcrSpaceResponse struct {
	ParsedResults []struct {
		TextOverlay struct {
			Lines []struct {
				LineText string `json:"LineText"`
				Words    []struct {
					WordText string  `json:"WordText"`
					Left     float64 `json:"Left"`
					Top      float64 `json:"Top"`
					Height   float64 `json:"Height"`
					Width    float64 `json:"Width"`
				} `json:"Words"`
				MaxHeight float64 `json:"MaxHeight"`
				MinTop    float64 `json:"MinTop"`
			} `json:"Lines"`
			HasOverlay bool   `json:"HasOverlay"`
			Message    string `json:"Message"`
		} `json:"TextOverlay"`
		TextOrientation   string `json:"TextOrientation"`
		FileParseExitCode int    `json:"FileParseExitCode"`
		ParsedText        string `json:"ParsedText"`
		ErrorMessage      string `json:"ErrorMessage"`
		ErrorDetails      string `json:"ErrorDetails"`
	} `json:"ParsedResults"`
	OCRExitCode                  int    `json:"OCRExitCode"`
	IsErroredOnProcessing        bool   `json:"IsErroredOnProcessing"`
	ProcessingTimeInMilliseconds string `json:"ProcessingTimeInMilliseconds"`
	SearchablePDFURL             string `json:"SearchablePDFURL"`
}

func (o OcrSpaceOcr) OCR(ctx context.Context, reder io.Reader, params map[string]string, filename string) (or OcrSpaceResponse, err error) {
	buf := &bytes.Buffer{}

	multip := multipart.NewWriter(buf)
	defer multip.Close()
	filew, err := multip.CreateFormFile("file", filename)
	if err != nil {
		return or, err
	}
	if _, err := io.Copy(filew, reder); err != nil {
		return or, err
	}

	for key, value := range params {
		if err := multip.WriteField(key, value); err != nil {
			return or, err
		}
	}

	multip.Close()

	request, err := http.NewRequestWithContext(ctx, "POST", o.ServerURL, buf)
	if err != nil {
		return or, err
	}
	request.Header.Set("Content-Type", multip.FormDataContentType())
	request.Header.Set("apikey", o.Token)

	response, err := o.Client.Do(request)
	if err != nil {
		return or, err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&or)
	return
}

func (o OcrSpaceOcr) Text(ctx context.Context, r io.Reader) (string, error) {
	response, err := o.OCR(ctx, r, o.Body, "image.png")
	if err != nil {
		return "", err
	}

	rj, _ := json.MarshalIndent(response, "", " ")
	fmt.Println(string(rj))

	if response.OCRExitCode != 1 {
		return "", errors.New("unexcepted ocr error")
	}

	if len(response.ParsedResults) == 0 {
		return "", nil
	}

	if len(response.ParsedResults[0].TextOverlay.Lines) == 0 {
		return "", nil
	}

	//return response.ParsedResults[0].TextOverlay.Lines[0].LineText, nil
	return strings.TrimSpace(response.ParsedResults[0].ParsedText), nil

}
