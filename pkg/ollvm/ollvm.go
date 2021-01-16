package ollvm

import (
	"bytes"
	"os"
	"os/exec"
	"path"
)

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

// Clang - Holds an instances of a clang install
type Clang struct {
	ClangRootDir string
	ClangExe     string
}

// InitClang - Initalize a Clang struct
func InitClang(clangDir string) (*Clang, error) {
	clang := &Clang{
		ClangRootDir: clangDir,
		ClangExe:     path.Join(clangDir, "bin", "clang.exe"),
	}
	if _, err := os.Stat(clang.ClangRootDir); os.IsNotExist(err) {
		return nil, err
	}
	if _, err := os.Stat(clang.ClangExe); os.IsNotExist(err) {
		return nil, err
	}
	return clang, nil
}

// Version - Get clang version info
func (c *Clang) Version() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	stdout, stderr, err := c.clangCmd(cwd, os.Environ(), []string{"--version"})
	if err != nil {
		return string(stderr), err
	}
	return string(stdout), nil
}

// clangCmd - Execute a nim command
func (c *Clang) clangCmd(wd string, env []string, command []string) ([]byte, []byte, error) {
	cmd := exec.Command(c.ClangExe, command...)
	cmd.Dir = wd
	cmd.Env = env
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}
