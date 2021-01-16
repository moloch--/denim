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
	"math/rand"
	"path/filepath"
	"time"

	"github.com/moloch--/denim/pkg/assets"
	"github.com/moloch--/denim/pkg/build"
	"github.com/moloch--/denim/pkg/nim"
	"github.com/moloch--/denim/pkg/ollvm"
	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile a nim program",
	Long:  `Compile a nim program with obfuscator-llvm`,
	Run: func(cmd *cobra.Command, args []string) {
		if !preflight() {
			return
		}

		if len(args) < 1 {
			fmt.Printf(Warn + "Missing input files\n")
			return
		}

		obfArgs, err := getObfArgs(cmd)
		if err != nil {
			return
		}

		build.Compile(filepath.Base(args[0]), args, obfArgs)
	},
}

func preflight() bool {
	_, err := nim.Version()
	if err != nil {
		fmt.Printf(Warn + "Could not find nim on PATH\n")
		return false
	}
	_, err = ollvm.InitClang(assets.GetClangDir())
	if err != nil {
		fmt.Printf(Warn + "No obfuscator-llvm found, you probably need to run 'denim setup'\n")
		return false
	}
	return true
}

func getObfArgs(cmd *cobra.Command) (*ollvm.ObfArgs, error) {
	obfArgs := &ollvm.ObfArgs{}
	rand.Seed(time.Now().UTC().UnixNano())

	bcfEnabled, err := cmd.Flags().GetBool(bcfFlagStr)
	if err != nil {
		fmt.Printf(Warn+"Failed to parse --%s flag: %s\n", bcfFlagStr, err)
		return nil, err
	}
	if bcfEnabled {
		obfArgs.BCF = bcfEnabled

		bcfLoops, err := cmd.Flags().GetInt(bcfLoopFlagStr)
		if err != nil {
			fmt.Printf(Warn+"Failed to parse --%s flag: %s\n", bcfLoopFlagStr, err)
			return nil, err
		}
		if bcfLoops < 1 {
			bcfLoops = rand.Intn(4) + 1
		}
		obfArgs.BCFLoop = bcfLoops

		bcfProb, err := cmd.Flags().GetInt(bcfProbFlagStr)
		if err != nil {
			fmt.Printf(Warn+"Failed to parse --%s flag: %s\n", bcfProbFlagStr, err)
			return nil, err
		}
		if 100 < bcfProb {
			fmt.Printf(Info + "Max BCF probability set to 100\n")
			bcfProb = 100
		}
		obfArgs.BCFProb = bcfProb
	}

	subEnabled, err := cmd.Flags().GetBool(subFlagStr)
	if err != nil {
		fmt.Printf(Warn+"Failed to parse --%s flag: %s\n", subFlagStr, err)
		return nil, err
	}
	if subEnabled {
		obfArgs.Sub = subEnabled

		subLoops, err := cmd.Flags().GetInt(subLoopFlagStr)
		if err != nil {
			fmt.Printf(Warn+"Failed to parse --%s flag: %s\n", subLoopFlagStr, err)
			return nil, err
		}
		if subLoops < 1 {
			subLoops = rand.Intn(2) + 1
		}
		obfArgs.SubLoop = subLoops
	}

	flattenEnabled, err := cmd.Flags().GetBool(flattenFlagStr)
	if err != nil {
		fmt.Printf(Warn+"Failed to parse --%s flag: %s\n", flattenFlagStr, err)
		return nil, err
	}
	if flattenEnabled {
		obfArgs.Flatten = flattenEnabled

		splits, err := cmd.Flags().GetInt(flattenSplitStr)
		if err != nil {
			fmt.Printf(Warn+"Failed to parse --%s flag: %s\n", flattenSplitStr, err)
			return nil, err
		}
		if splits < 1 {
			splits = rand.Intn(4) + 1
		}
		obfArgs.FlattenSplit = splits
	}

	seed, err := cmd.Flags().GetString(seedFlagStr)
	if err != nil {
		fmt.Printf(Warn+"Failed to parse --%s flag: %s\n", seedFlagStr, err)
		return nil, err
	}
	obfArgs.AESSeed = seed

	return obfArgs, nil
}
