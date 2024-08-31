package settings

import (
	"KamaiZen/logger"
	"encoding/json"
	"os"
)

type LSPSettings struct {
	KamailioSourcePath string          `json:"kamailioSourcePath"`
	RootDir            string          `json:"rootDir"`
	LogLevel           logger.LOGLEVEL `json:"logLevel"`
}

// NewLSPSettings creates and returns a new instance of LSPSettings.
// It initializes the settings with the given Kamailio source path, root directory, and log level.
//
// Parameters:
//
//	kamailioSourcePath string - The path to the Kamailio source code.
//	rootDir string - The root directory for the language server.
//	log_level logger.LOGLEVEL - The logging level for the language server.
//
// Returns:
//
//	LSPSettings - The initialized settings.
func NewLSPSettings(kamailioSourcePath string, rootDir string, log_level logger.LOGLEVEL) LSPSettings {
	return LSPSettings{
		KamailioSourcePath: kamailioSourcePath,
		RootDir:            rootDir,
		LogLevel:           log_level,
	}
}

type JSONSettingsReader struct{}

// ReadSettings reads the settings from a JSON file at the given filepath.
// It unmarshals the JSON data into an LSPSettings struct and returns it.
//
// Parameters:
//
//	filepath string - The path to the JSON file containing the settings.
//
// Returns:
//
//	LSPSettings - The settings read from the JSON file.
func (jsr *JSONSettingsReader) ReadSettings(filepath string) LSPSettings {
	data, err := os.ReadFile(filepath)
	if err != nil {
		logger.Errorf("Error reading file: %s", err)
	}
	var settings LSPSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		logger.Errorf("Error unmarshalling settings: %s", err)
	}
	return settings
}

const RPC_VERSION = "2.0"
const KAMAIZEN_VERSION = "0.0.1"
const MY_NAME = "KamaiZen"
