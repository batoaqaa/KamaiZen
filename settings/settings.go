package settings

type LSPSettings struct {
	kamailioSourcePath string
	rootDir            string
}

func NewLSPSettings(kamailioSourcePath string, rootDir string) LSPSettings {
	return LSPSettings{
		kamailioSourcePath: kamailioSourcePath,
		rootDir:            rootDir,
	}
}

func (s *LSPSettings) KamailioSourcePath() string {
	return s.kamailioSourcePath
}

func (s *LSPSettings) RootDir() string {
	return s.rootDir
}
