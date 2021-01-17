package cmd

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

	"github.com/moloch--/denim/pkg/assets"
	"github.com/moloch--/denim/pkg/nim"
	"github.com/moloch--/denim/pkg/ollvm"
	"github.com/spf13/cobra"
)

var (
	// Version - The semantic version of the program
	Version string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Print the version number of denim and exit`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Denim v%s\n\n", Version)

		nimVer, err := nim.Version()
		if err != nil {
			fmt.Printf(Warn + "Nim does not appear to be on your PATH!\n")
		} else {
			fmt.Printf(nimVer + "\n\n")
		}

		clang, err := ollvm.InitClang(assets.GetClangDir())
		if err != nil {
			fmt.Printf(Warn + "No clang, please run 'denim setup'")
		} else {
			clangVer, err := clang.Version()
			if err != nil {
				fmt.Printf(Warn+"%s\n", err)
			} else {
				fmt.Printf(clangVer + "\n")
			}
		}

	},
}
