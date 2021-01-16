package ollvm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/moloch--/denim/pkg/assets"
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
const (

	// MaxProb - Maxium probability
	MaxProb = 100
	// MaxBCFLoop - Max loop value
	MaxBCFLoop = 5
	// MaxSubLoop - Max loop value
	MaxSubLoop = 4
	// MaxSplit - Max split value
	MaxSplit = 5
)

// Clang - Holds an instances of a clang install
type Clang struct {
	ClangRootDir string
	ClangBinDir  string
	ClangExe     string
}

// ObfArgs - Build options
type ObfArgs struct {

	// Compile Options
	Static bool     `json:"static,omitempty"`
	Link   []string `json:"link,omitempty"`

	// Bogus control flow
	BCF     bool `json:"bcf"`
	BCFProb int  `json:"bcf_prob"`
	BCFLoop int  `json:"bcf_loop"`

	// Instructions substitution
	Sub     bool `json:"sub"`
	SubLoop int  `json:"sub_loop"`

	// Control flow flattening
	Flatten      bool `json:"flaten"`
	FlattenSplit int  `json:"flatten_split"`

	AESSeed string `json:"aes_seed"`
}

// InitClang - Initalize a Clang struct
func InitClang(clangDir string) (*Clang, error) {
	clang := &Clang{
		ClangRootDir: clangDir,
		ClangBinDir:  path.Join(clangDir, "bin"),
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
	mingw := filepath.Join(assets.GetMingwDir(), "bin")
	env := []string{
		fmt.Sprintf("PATH=%s;%s", c.ClangBinDir, mingw),
	}
	stdout, stderr, err := c.clangCmd(cwd, env, []string{"--version"})
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

func (c *Clang) verifyObfArgs(obfArgs *ObfArgs) error {
	if obfArgs.BCF && MaxProb < obfArgs.BCFProb {
		return fmt.Errorf("BFC probability cannot exceed %d", MaxProb)
	}
	if obfArgs.BCF && MaxBCFLoop < obfArgs.BCFLoop {
		return fmt.Errorf("BCF loop cannot exceed %d", MaxBCFLoop)
	}
	if obfArgs.Sub && MaxSubLoop < obfArgs.SubLoop {
		return fmt.Errorf("Substitution loop cannot exceed %d", MaxSubLoop)
	}
	if obfArgs.Flatten && MaxSplit < obfArgs.FlattenSplit {
		return fmt.Errorf("Flatten split cannot exceed %d", MaxSplit)
	}
	return nil
}
