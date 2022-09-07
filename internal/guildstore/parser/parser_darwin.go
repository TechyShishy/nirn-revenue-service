package parser

import (
	"log"
	"os"
	"path/filepath"
)

func init() {
	homePath, ok := os.LookupEnv("HOME")
	if !ok {
		log.Print("could not lookup $HOME")
		return
	}
	DefaultGSDataFileGlob = filepath.Join(
		homePath,
		defaultSavedVariablesPathBase,
		defaultGSDataFileGlobBase,
	)
}
