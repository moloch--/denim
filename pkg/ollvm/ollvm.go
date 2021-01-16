package ollvm

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
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
	mingwBinDir := filepath.Join(assets.GetMingwDir(), "bin")
	env := []string{
		fmt.Sprintf("PATH=%s;%s", c.ClangBinDir, mingwBinDir),
	}
	stdout, stderr, err := c.clangCmd(cwd, env, []string{"--version"})
	if err != nil {
		return string(stderr), err
	}
	return string(stdout), nil
}

// ObfCompile - Compile obfuscated C code
func (c *Clang) ObfCompile(wd string, args []string, obfArgs *ObfArgs) ([]byte, []byte, error) {
	err := c.verifyObfArgs(obfArgs)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	mingwBinDir := filepath.Join(assets.GetMingwDir(), "bin")
	env := []string{
		fmt.Sprintf("PATH=%s;%s", c.ClangBinDir, mingwBinDir),
	}
	command := c.getCmdObfArgs(obfArgs)
	command = append(command, args...)
	return c.clangCmd(wd, env, command)
}

// Compile - Compile C code (no obfuscation)
func (c *Clang) Compile(wd string, args []string) ([]byte, []byte, error) {
	mingwBinDir := filepath.Join(assets.GetMingwDir(), "bin")
	env := []string{
		fmt.Sprintf("PATH=%s;%s", c.ClangBinDir, mingwBinDir),
	}
	return c.clangCmd(wd, env, args)
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

func (c *Clang) getCmdObfArgs(obfArgs *ObfArgs) []string {
	cmdArgs := []string{}

	if obfArgs.BCF {
		cmdArgs = append(cmdArgs, []string{"-mllvm", "-bcf"}...)
		bcfProb := fmt.Sprintf("-bcf_prob=%d", getIntArg(obfArgs.BCFProb))
		cmdArgs = append(cmdArgs, []string{"-mllvm", bcfProb}...)
		bcfLoop := fmt.Sprintf("-bcf_loop=%d", getIntArg(obfArgs.BCFLoop))
		cmdArgs = append(cmdArgs, []string{"-mllvm", bcfLoop}...)
	}
	if obfArgs.Sub {
		cmdArgs = append(cmdArgs, []string{"-mllvm", "-sub"}...)
		subLoop := fmt.Sprintf("-sub_loop=%d", getIntArg(obfArgs.SubLoop))
		cmdArgs = append(cmdArgs, []string{"-mllvm", subLoop}...)
	}
	if obfArgs.AESSeed == "" {
		obfArgs.AESSeed = randomSeed()
	}
	digest := sha256.New()
	digest.Write([]byte(obfArgs.AESSeed))
	aesSeed := fmt.Sprintf("-aesSeed=%x", digest.Sum(nil)[:16])
	cmdArgs = append(cmdArgs, []string{"-mllvm", aesSeed}...)
	return cmdArgs
}

func getIntArg(x int) int {
	if x < 1 {
		return 1
	}
	return x
}

func randomSeed() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	digest := sha256.New()
	digest.Write(buf)
	return fmt.Sprintf("%x", digest.Sum(nil))
}
