package plugins

import (
	"MediaDispenserGo/middlewares"
	"bytes"
	"errors"
	"slices"
)

// FileProcessor je interface pro zpracování souborů pluginem.
type FileProcessor interface {
	Process(file bytes.Buffer, mimeType string, fileName string, fileExt string) (processedFile []byte, newMimeType string, newName string, newExtension string, err error)
	GetName() string
}

type RegistryRecord struct {
	Constructor       func() FileProcessor
	AllowedOperations []middlewares.Operation
}

// Mapování pluginů podle jejich názvu
var registeredPlugins = map[string]RegistryRecord{}

// RegisterPlugin umožňuje registrovat nový plugin.
func RegisterPlugin(name string, constructor func() FileProcessor, allowedOperations ...middlewares.Operation) {
	registeredPlugins[name] = RegistryRecord{
		constructor,
		allowedOperations,
	}
}

// GetPluginInstance vrací instanci pluginu podle jeho názvu.
func GetPluginInstance(name string, op middlewares.Operation) (FileProcessor, error) {
	record, exists := registeredPlugins[name]
	if !exists {
		return nil, errors.New("plugin not found")
	}

	if !slices.Contains(record.AllowedOperations, op) {
		return nil, errors.New("plugin are not allowed in this operation")
	}

	return record.Constructor(), nil
}

// Dostupné pluginy zaregistrujeme při inicializaci balíčku
func init() {
	RegisterPlugin("WebPConverter", func() FileProcessor { return &WebPConverter{} }, middlewares.OperationUpload)
	RegisterPlugin("KSUIDRenamer", func() FileProcessor { return &KSUIDRenamer{} }, middlewares.OperationUpload)
}
