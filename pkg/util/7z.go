package util

import (
	"fmt"
	"os/exec"
)

// Extract7z - Extract a 7z archive using the console 7z util
func Extract7z(sevenZipExe string, archive string, dest string) error {
	cmd := exec.Command(sevenZipExe, []string{"x", fmt.Sprintf("-o%s", dest), archive}...)
	return cmd.Run()
}
