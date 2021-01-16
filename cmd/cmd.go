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
	"os"

	"github.com/spf13/cobra"
)

const (
	// ANSI Colors
	normal    = "\033[0m"
	black     = "\033[30m"
	red       = "\033[31m"
	green     = "\033[32m"
	orange    = "\033[33m"
	blue      = "\033[34m"
	purple    = "\033[35m"
	cyan      = "\033[36m"
	gray      = "\033[37m"
	bold      = "\033[1m"
	clearln   = "\r\x1b[2K"
	upN       = "\033[%dA"
	downN     = "\033[%dB"
	underline = "\033[4m"

	// Info - Display colorful information
	Info = bold + cyan + "[*] " + normal
	// Warn - Warn a user
	Warn = bold + red + "[!] " + normal
	// Debug - Display debug information
	Debug = bold + purple + "[-] " + normal
	// Woot - Display success
	Woot = bold + green + "[$] " + normal

	// Setup - Standard Flags
	timeoutFlagStr           = "timeout"
	skipTLSValidationFlagStr = "skip-tls-validation"
	proxyFlagStr             = "proxy"

	// Compile - Standard Flags
	outputFlagStr = "output"

	// Compile - Obfuscation Flags
	bcfFlagStr      = "bcf"
	bcfLoopFlagStr  = "bcf-loop"
	bcfProbFlagStr  = "bcf-probability"
	subFlagStr      = "sub"
	subLoopFlagStr  = "sub-loop"
	flattenFlagStr  = "flatten"
	flattenSplitStr = "flatten-split"
	seedFlagStr     = "seed"
)

var rootCmd = &cobra.Command{
	Use:   "denim",
	Short: "Automated compiler obfuscation for nim",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {

	// Version
	rootCmd.AddCommand(versionCmd)

	// Setup options
	setupCmd.Flags().BoolP(skipTLSValidationFlagStr, "V", false, "Skip TLS certificate validation")
	setupCmd.Flags().StringP(proxyFlagStr, "H", "", "Specify HTTP(S) proxy URL (e.g. http://localhost:8080)")
	setupCmd.Flags().IntP(timeoutFlagStr, "T", 3600, "HTTPS request/connection timeout (default: 1hr)")
	rootCmd.AddCommand(setupCmd)

	// Compile - Obfuscator options
	compileCmd.Flags().BoolP(bcfFlagStr, "b", true, "Enable bogus control flow")
	compileCmd.Flags().IntP(bcfLoopFlagStr, "C", 0, "Number of bogus control flow passes (0 = random)")
	compileCmd.Flags().IntP(bcfProbFlagStr, "F", 100, "Probability a basic bloc will be obfuscated")
	compileCmd.Flags().BoolP(subFlagStr, "s", true, "Enable instruction substitution")
	compileCmd.Flags().IntP(subLoopFlagStr, "U", 0, "Number of instruction substitution passes (0 = random)")
	compileCmd.Flags().BoolP(flattenFlagStr, "f", true, "Enable control flow flattening")
	compileCmd.Flags().IntP(flattenSplitStr, "L", 0, "Splits applied to each block (0 = random)")
	compileCmd.Flags().StringP(seedFlagStr, "r", "", "PRNG obfuscation seed (default is random)")

	// Compile - Standard options
	compileCmd.Flags().StringP(outputFlagStr, "o", "", "output file")
	rootCmd.AddCommand(compileCmd)

}

// Execute - Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
