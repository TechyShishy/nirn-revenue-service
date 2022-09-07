package parser

import (
	"log"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func init() {
	documentsPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, 0)
	if err != nil {
		log.Print(err)
		return
	}
	DefaultGSDataFileGlob = filepath.Join(
		documentsPath,
		defaultSavedVariablesPathBase,
		defaultGSDataFileGlobBase,
	)
}
