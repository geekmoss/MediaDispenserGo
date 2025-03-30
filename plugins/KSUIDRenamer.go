package plugins

import (
	"bytes"
	"github.com/segmentio/ksuid"
)

type KSUIDRenamer struct{}

func (k *KSUIDRenamer) Process(file bytes.Buffer, mimeType string, _ string, fileExt string) ([]byte, string, string, string, error) {
	// Generujeme nový název na základě KSUID
	newID := ksuid.New()
	newFilename := newID.String()

	// Vracíme původní obsah a nový název (bez přípony)
	return file.Bytes(), mimeType, newFilename, fileExt, nil
}

func (k *KSUIDRenamer) GetName() string {
	return "KSUIDRenamer"
}
