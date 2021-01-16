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
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cheggaaa/pb/v3"
	"github.com/moloch--/denim/pkg/assets"
	"github.com/moloch--/denim/pkg/nim"
	"github.com/moloch--/denim/pkg/util"
	"github.com/spf13/cobra"
)

var (
	// ObfuscatorLLVMURL - URL to a O-LLVM Github repo
	ObfuscatorLLVMURL string

	// NimURL - URL to the nim package
	NimURL string
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup denim",
	Long:  `Download obfuscator-llvm and nim tool chains`,
	Run: func(cmd *cobra.Command, args []string) {
		setup(cmd, args)
	},
}

func setup(cmd *cobra.Command, args []string) {
	denimDir := assets.GetRootDir()

	_, err := nim.Version()
	if err != nil {
		fmt.Printf(Warn + "Nim does not appear to be on your PATH!")
	}

	client := initHTTPClient(cmd)
	if client == nil {
		return
	}

	fmt.Println("Downloading obfuscator llvm ...")
	llvmTar := filepath.Join(denimDir, "ollvm.tar.gz")
	if _, err := os.Stat(llvmTar); !os.IsNotExist(err) {
		os.Remove(llvmTar)
	}
	err = downloadAsset(client, ObfuscatorLLVMURL, llvmTar)
	if err != nil {
		fmt.Printf(Warn+"Download failed %s\n", err)
	}

	fmt.Println("Extracting obfuscator-llvm, please wait ...")
	unpackDir := filepath.Join(denimDir, "ollvm")
	if _, err := os.Stat(unpackDir); !os.IsNotExist(err) {
		os.RemoveAll(unpackDir)
	}
	tarReader, err := os.Open(llvmTar)
	if err != nil {
		fmt.Printf(Warn+"Failed to read %s", err)
		return
	}
	err = util.Untar(unpackDir, tarReader)
	if err != nil {
		fmt.Printf(Warn+"Failed to extract obfuscator-llvm %s", err)
		return
	}
	tarReader.Close()
	os.Remove(llvmTar)
}

func initHTTPClient(cmd *cobra.Command) *http.Client {
	timeoutSeconds, err := cmd.Flags().GetInt("timeout")
	timeout := time.Duration(timeoutSeconds * int(time.Second))
	if err != nil {
		fmt.Printf("Failed to parse --timeout flag: %s\n", err)
		return nil
	}

	skipTLSValidation, err := cmd.Flags().GetBool("skip-tls-validation")
	if err != nil {
		fmt.Printf("Failed to parse --skip-tls-validation flag: %s\n", err)
		return nil
	}
	if skipTLSValidation {
		fmt.Println()
		fmt.Println(Warn + "You're trying to download the compilers over an insecure connection, this is a bad idea!")
		confirm := false
		prompt := &survey.Confirm{Message: "Continue?"}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			return nil
		}
		confirm = false
		prompt = &survey.Confirm{Message: "Seriously?"}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			return nil
		}
	}

	proxy, err := cmd.Flags().GetString("proxy")
	if err != nil {
		fmt.Printf("Failed to parse --proxy flag: %s\n", err)
		return nil
	}
	var proxyURL *url.URL = nil
	if proxy != "" {
		proxyURL, err = url.Parse(proxy)
		if err != nil {
			fmt.Printf(Warn+"%s", err)
			return nil
		}
	}

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: timeout,
			}).Dial,
			TLSHandshakeTimeout: timeout,
			Proxy:               http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipTLSValidation,
			},
		},
	}

	return client
}

func downloadAsset(client *http.Client, assetURL string, saveTo string) error {
	writer, err := os.Create(saveTo)
	if err != nil {
		return err
	}
	resp, err := client.Get(assetURL)
	if err != nil {
		return err
	}
	bar := pb.Full.Start64(resp.ContentLength)
	barReader := bar.NewProxyReader(resp.Body)
	io.Copy(writer, barReader)
	bar.Finish()
	fmt.Printf(upN, 1)
	fmt.Printf(clearln + "\r")
	return nil
}
