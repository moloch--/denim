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
	"fmt"
	"os"
	"path/filepath"

	"github.com/moloch--/denim/pkg/assets"
	"github.com/moloch--/denim/pkg/nim"
	"github.com/moloch--/denim/pkg/ollvm"
)

// Compile a nim program with Obfuscator-LLVM
func Compile(project string, nimFiles []string, obfArgs *ollvm.ObfArgs) error {
	clang, err := ollvm.InitClang(assets.GetClangDir())
	if err != nil {
		return err
	}
	nimCache, err := compileNimCode(project, nimFiles, clang)
	if err != nil {
		return err
	}

	parseProjectJson(nimCache)

	return nil
}

//  nim compile --genScript --compileOnly --cc=clang --clang.exe:PATH --nimcache:PATH helloworld.nim
func compileNimCode(project string, nimFiles []string, clang *ollvm.Clang) (string, error) {
	nimCache := filepath.Join(assets.GetNimCacheRoot(), project)
	args := []string{"--genScript", "--compileOnly", "--cc:clang"}
	args = append(args, fmt.Sprintf("--clang.exe=%s", clang.ClangExe))
	args = append(args, fmt.Sprintf("--nimcache:%s", nimCache))
	args = append(args, nimFiles...)

	fmt.Println("------------[nim]------------")
	fmt.Printf(" > nim compile %v\n\n", args)

	workDir, _ := os.Getwd()
	stdout, stderr, err := nim.Compile(workDir, os.Environ(), args)

	fmt.Printf(string(stdout))
	fmt.Printf(string(stderr))
	fmt.Println("-----------------------------")

	return nimCache, err
}

func parseProjectJson(nimCache string) {

}
