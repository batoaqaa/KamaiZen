package settings

import "KamaiZen/logger"

type LSPSettings struct {
	kamailioSourcePath string
	rootDir            string
	logLevel           logger.LOGLEVEL
}

func NewLSPSettings(kamailioSourcePath string, rootDir string, log_level logger.LOGLEVEL) LSPSettings {
	return LSPSettings{
		kamailioSourcePath: kamailioSourcePath,
		rootDir:            rootDir,
		logLevel:           log_level,
	}
}

func (s *LSPSettings) KamailioSourcePath() string {
	return s.kamailioSourcePath
}

func (s *LSPSettings) RootDir() string {
	return s.rootDir
}

func (s *LSPSettings) LogLevel() logger.LOGLEVEL {
	return s.logLevel
}

const RPC_VERSION = "2.0"
const KAMAIZEN_VERSION = "0.0.1"
const MY_NAME = "KamaiZen"
