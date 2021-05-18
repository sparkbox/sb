// Package to house shared code between cmds.
package util

import (
	"io/ioutil"

	homedir "github.com/mitchellh/go-homedir"
)

// Returns the token as a string or an error
func ReturnToken() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	token, err := ioutil.ReadFile(home + "/.sb")

	if err != nil {
		return "", err
	}

	return string(token), nil
}
