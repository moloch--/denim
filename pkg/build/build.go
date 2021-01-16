package build

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/moloch--/denim/pkg/assets"
	"github.com/moloch--/denim/pkg/nim"
	"github.com/moloch--/denim/pkg/ollvm"
)

// Build - Denim build args
type Build struct {
	Name     string
	NimFiles []string

	UserCodeOnly bool
}

// Compile a nim program with Obfuscator-LLVM
func Compile(build *Build, obfArgs *ollvm.ObfArgs) error {
	clang, err := ollvm.InitClang(assets.GetClangDir())
	if err != nil {
		return err
	}

	// Compile Nim
	nimCache, err := compileNimCode(build.Name, build.NimFiles, clang)
	if err != nil {
		return err
	}
	nimProject, err := parseProjectJSON(nimCache)
	if err != nil {
		return err
	}

	// Compile C
	for _, step := range nimProject.Compile {
		if len(step) != 2 {
			return fmt.Errorf("Malformed step: %v", step)
		}
		cFile := filepath.Base(step[0])
		compileCmd := strings.Fields(step[1])
		if compileCmd[0] == "clang" || compileCmd[0] == "clang.exe" {
			compileCmd = compileCmd[1:]
		}
		var err error
		if strings.HasPrefix(cFile, "@") || !build.UserCodeOnly {
			_, _, err = clang.ObfCompile(nimCache, compileCmd, obfArgs)
		} else {
			_, _, err = clang.Compile(nimCache, compileCmd)
		}
		if err != nil {
			return err
		}
	}
	linker := strings.Fields(nimProject.LinkCmd)
	if linker[0] == "clang" || linker[0] == "clang.exe" {
		linker = linker[1:]
	}
	_, _, err = clang.Compile(nimCache, linker)
	if err != nil {
		return err
	}

	return nil
}

// nim compile --genScript --compileOnly --cc=clang --clang.exe:PATH --nimcache:PATH helloworld.nim
func compileNimCode(project string, nimFiles []string, clang *ollvm.Clang) (string, error) {
	nimCache := filepath.Join(assets.GetNimCacheRoot(), project)
	args := []string{"--genScript", "--compileOnly", "--cc:clang"}
	args = append(args, fmt.Sprintf("--clang.exe=%s", clang.ClangExe))
	args = append(args, fmt.Sprintf("--nimcache:%s", nimCache))
	args = append(args, nimFiles...)
	workDir, _ := os.Getwd()
	_, _, err := nim.Compile(workDir, os.Environ(), args)
	return nimCache, err
}

func parseProjectJSON(nimCache string) (*nim.Project, error) {
	entries, err := ioutil.ReadDir(nimCache)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		entryPath := filepath.Join(nimCache, entry.Name())
		if strings.HasSuffix(entryPath, ".json") {
			data, err := ioutil.ReadFile(entryPath)
			if err != nil {
				return nil, err
			}
			project := &nim.Project{}
			err = json.Unmarshal(data, project)
			return project, err
		}
	}
	return nil, nil
}
