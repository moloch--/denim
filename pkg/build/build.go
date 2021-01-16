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

	"github.com/moloch--/denim/pkg/assets"
	"github.com/moloch--/denim/pkg/nim"
	"github.com/moloch--/denim/pkg/ollvm"
)

//  nim compile --genScript --compileOnly --cc=clang --clang.exe:PATH --nimcache:PATH helloworld.nim

// Compile a nim program with Obfuscator-LLVM
func Compile(nimFiles []string, obfArgs *ollvm.ObfArgs) error {
	clang, err := ollvm.InitClang(assets.GetClangDir())
	if err != nil {
		return err
	}
	compileNimCode(nimFiles, clang)

	return nil
}

func compileNimCode(nimFiles []string, clang *ollvm.Clang) error {
	args := []string{"--genScript", "--compileOnly", "--cc:clang"}
	args = append(args, fmt.Sprintf("--clang.exe=%s", clang.ClangExe))
	args = append(args, fmt.Sprintf("--nimcache:%s", assets.GetNimCache()))
	args = append(args, nimFiles...)

	workDir, _ := os.Getwd()
	stdout, stderr, err := nim.Compile(workDir, os.Environ(), args)

	fmt.Println("------------[nim]------------")
	fmt.Printf(string(stdout))
	fmt.Printf(string(stderr))
	fmt.Println("-----------------------------")

	return err
}
