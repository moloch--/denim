package assets

/*
	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.
	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.
	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
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

// GetClangDir - Get the clang root directory
func GetClangDir() string {
	rootDir := GetRootDir()
	return filepath.Join(rootDir, "ollvm", "build")
}

// GetNimCache - Get the clang root directory
func GetNimCache() string {
	rootDir := GetRootDir()
	nimcache := filepath.Join(rootDir, "nimcache")
	if _, err := os.Stat(nimcache); os.IsNotExist(err) {
		err = os.MkdirAll(nimcache, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nimcache
}
