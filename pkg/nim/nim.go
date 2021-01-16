package nim

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
	"bytes"
	"os"
	"os/exec"
)

// nimCmd - Execute a nim command
func nimCmd(wd string, env []string, command []string) ([]byte, []byte, error) {
	cmd := exec.Command(Nim, command...)
	cmd.Dir = wd
	cmd.Env = env
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

// Version - Get nim version output
func Version() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	stdout, stderr, err := nimCmd(cwd, []string{}, []string{"--version"})
	if err != nil {
		return string(stderr), err
	}
	return string(stdout), nil
}

// Compile - Nim compiler command
func Compile(workDir string, env []string, args []string) ([]byte, []byte, error) {
	cli := []string{"compile"}
	cli = append(cli, args...)
	return nimCmd(workDir, env, cli)
}
