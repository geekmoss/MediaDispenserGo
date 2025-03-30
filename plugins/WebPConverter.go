package plugins

import (
	"bytes"
	"github.com/chai2010/webp"
	"image"
	"image/jpeg"
	"image/png"
)

type WebPConverter struct{}

func (w *WebPConverter) Process(file bytes.Buffer, mimeType string, fileName string, fileExt string) ([]byte, string, string, string, error) {
	var img image.Image
	var err error

	// Dekódování obrázku podle MIME typu
	switch mimeType {
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(file.Bytes()))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(file.Bytes()))
	default:
		return file.Bytes(), mimeType, fileName, ".webp", nil
	}

	if err != nil {
		return nil, "", "", "", err
	}

	// Převedení obrázku na WebP
	var buf bytes.Buffer
	err = webp.Encode(&buf, img, &webp.Options{Lossless: true})
	if err != nil {
		return nil, "", "", "", err
	}

	// Vrácení konvertovaného souboru
	return buf.Bytes(), "image/webp", fileName, ".webp", nil
}

func (w *WebPConverter) GetName() string {
	return "WebPConverter"
}
