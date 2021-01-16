package cmd

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
	client := initHTTPClient(cmd)
	if client == nil {
		return
	}

	fmt.Println("Downloading nim ...")
	err := downloadAsset(client, NimURL, denimDir)
	if err != nil {
		fmt.Printf(Warn+"Download failed %s\n", err)
	}

	fmt.Println("Downloading obfuscator llvm ...")
	err = downloadAsset(client, ObfuscatorLLVMURL, denimDir)
	if err != nil {
		fmt.Printf(Warn+"Download failed %s\n", err)
	}
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
	downloadURL, err := url.Parse(assetURL)
	if err != nil {
		return err
	}
	assetFileName := filepath.Base(downloadURL.Path)
	writer, err := os.Create(filepath.Join(saveTo, assetFileName))
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
	return nil
}
