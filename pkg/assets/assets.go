package assets

import (
	"log"
	"os"
	"os/user"
	"path"
)

const (
	// DenimRootDirName - Directory storing all of the client configs/logs
	DenimRootDirName = ".denim"
)

// GetRootDir - Get the denim root directory
func GetRootDir() string {
	user, _ := user.Current()
	dir := path.Join(user.HomeDir, DenimRootDirName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir
}
