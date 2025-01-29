package setup

import (
	"os"

	"github.com/mainanick/shellby/internal/constants"
)

// CreateDir creates a directory if it does not exist
func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateAllDir() error {
	for _, path := range []string{
		constants.PATH,
		constants.TRAEFIK_PATH,
		constants.DYNAMIC_TRAEFIK_PATH,
		constants.LOG_PATH,
		constants.SSH_PATH} {
		if err := CreateDir(path); err != nil {
			return err
		}
	}
	return nil
}

func Initialize() {
	// Create all directories
	if err := CreateAllDir(); err != nil {
		panic(err)
	}

	// Setup Traefik

}
