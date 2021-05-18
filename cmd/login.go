package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"sb/util"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login via Slack",
	Long: `Login via Slack which returns your slackID and token.
You will need to paste the token and ID back to sb to where it will save
the value to ~/.sb`,
	Run: func(cmd *cobra.Command, args []string) {
		showToken, error := cmd.Flags().GetBool("token")

		if error != nil {
			log.Fatal("Error parsing token flag: ", error)
		}

		if showToken {
			token, err := util.ReturnToken()

			if err != nil {
				log.Fatal("failed to return token: ", err)
			}

			fmt.Println(token)
		} else {
			signInWithSlack()
		}
	},
}

func saveToken(token []byte) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(home+"/.sb", token, 0660)

	if err != nil {
		return err
	}

	return nil
}

func signInWithSlack() {
	var token []byte
	const url = `https://slack.com/oauth/v2/authorize?client_id=2160869413.2032676856630&user_scope=identity.basic`
	browser.OpenURL(url)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Paste in the returned value after authorizing with Slack: ")
	scanner.Scan()
	text := scanner.Bytes()
	if len(text) != 0 {
		token = text
	}

	err := saveToken(token)

	if err != nil {
		log.Fatal("Error saving token: ", err)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().BoolP("token", "t", false, "View saved token")
}
