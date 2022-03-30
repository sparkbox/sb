package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

func getURL() (string, error) {
	os := runtime.GOOS
	var url = "https://raw.githubusercontent.com/sparkbox/standard/main/security/security_policy_compliance/"

	switch os {
	case "darwin":
		return url + "policy-checks.sh", nil
	case "linux":
		return url + "linux-policy-check.sh", nil
	default:
		return "", errors.New("No security script for OS " + os)
	}
}

var securityCmd = &cobra.Command{
	Use:   "security",
	Short: "Run Security Policy Checks",
	Long:  `Run the Sparkbox security policy checks. Further details can be found at https://github.com/sparkbox/standard/tree/main/security/security_policy_compliance`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := getURL()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Running Security Policy Checks")
		fmt.Println("============================")
		fmt.Println("Getting sudo access...")

		response, err := http.Get(url)
		if err != nil {
			log.Fatalln("Failed to download script ", err)
		}

		defer response.Body.Close()

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalln("Failed to read script file ", err)
		}

		scriptFile, err := os.CreateTemp("", "security-policy-script")
		if err != nil {
			log.Fatal("Failed to create temp file ", err)
		}

		_, err = scriptFile.Write(bytes)
		if err != nil {
			log.Fatal("Failed to write script file: ", err)
		}

		err = scriptFile.Close()
		if err != nil {
			log.Fatal("Failed to close file: ", err)
		}

		err = os.Chmod(scriptFile.Name(), 0777)
		if err != nil {
			log.Fatal("Failed to make script executable ", err)
		}

		out, err := exec.Command(scriptFile.Name()).Output()
		if err != nil {
			log.Fatal("Failed to execute script", err)
		}

		fmt.Println(string(out))

		err = os.Remove(scriptFile.Name())
		if err != nil {
			log.Fatal("Failed to clean up after myself ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(securityCmd)
}
