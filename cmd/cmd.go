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

	// Setup
	// Proxy options
	setupCmd.Flags().BoolP("skip-tls-validation", "V", false, "Skip TLS certificate validation")
	setupCmd.Flags().StringP("proxy", "H", "", "Specify HTTP(S) proxy URL (e.g. http://localhost:8080)")
	setupCmd.Flags().IntP("timeout", "T", 30, "HTTPS request/connection timeout")
	rootCmd.AddCommand(setupCmd)

}

// Execute - Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
