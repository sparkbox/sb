package cmd

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sb/util"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type sshLogin struct {
	Token string
}

type Cert struct {
	Certificate string
	Key         string
}

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Generate a new SSH certificate",
	Long: `This generates a new SSH cert by talking to slackd via HTTP.
sb adds the returned Cert + Private key to your local ssh-agent.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")

		if err != nil {
			log.Fatal("Error parsing url flag: ", err)
		}

		print, err := cmd.Flags().GetBool("print")

		if err != nil {
			log.Fatal("Error parsing print flag: ", err)
		}

		write, err := cmd.Flags().GetBool("write")

		if err != nil {
			log.Fatal("Error parsing write flag: ", err)
		}

		cert, err := getCert(url)

		if err != nil {
			log.Fatal("Error getting cert: ", err)
		}

		if cert.Key == "" {
			log.Fatal("Cert doesn't have a Key.")
		}

		if print {
			printCertAndKey(cert)
		}

		if write {
			writeToFile(cert)
		}

		if !write && !print {
			sshCert, key, err := parseCert(cert)
			if err != nil {
				log.Fatal("Error parsing Certificate: ", err)
			}

			addToAgent(sshCert, key)
		}
	},
}

func writeToFile(cert Cert) {
	err := ioutil.WriteFile("key", []byte(cert.Key), 0600)
	if err != nil {
		fmt.Println("Error writing key file: ", err)
	}

	err = ioutil.WriteFile("key-cert.pub", []byte(cert.Certificate), 0600)
	if err != nil {
		fmt.Println("Error writing key file: ", err)
	}
}

func printCertAndKey(cert Cert) {
	fmt.Println(`ssh-agent expects a private key to go along with a certificate.
To add a certificate to ssh-agent manually you specify the private key (bar)
and it will automatically detect the certificate (bar-cert.pub) by convention.`)
	fmt.Println("Certificate: ")
	fmt.Println(cert.Certificate)
	fmt.Println("Key: ")
	fmt.Println(cert.Key)
}

func parseCert(cert Cert) (sshCert *ssh.Certificate, key *ecdsa.PrivateKey, error error) {
	block, _ := pem.Decode([]byte(cert.Key))

	key, err := x509.ParseECPrivateKey(block.Bytes)

	if err != nil {
		return sshCert, key, err
	}

	pub, _, _, _, err := ssh.ParseAuthorizedKey([]byte(cert.Certificate))

	if err != nil {
		return sshCert, key, err
	}

	return pub.(*ssh.Certificate), key, nil
}

func addToAgent(cert *ssh.Certificate, key *ecdsa.PrivateKey) {
	// validBefore - now === ssh-agent lifetime
	diff := time.Unix(int64(cert.ValidBefore), 0).Sub(time.Now())
	// need seconds for lifetime to pass to Agent
	lifetime := int(diff.Seconds())

	con, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))

	if err != nil {
		fmt.Println("can't connect to SSH agent: ", err)
	}

	sshAgent := agent.NewClient(con)

	err = sshAgent.Add(agent.AddedKey{
		PrivateKey:   key,
		Certificate:  cert,
		LifetimeSecs: uint32(lifetime),
	})

	if err != nil {
		log.Fatal("Failed to add to ssh-agent: ", err)
	}
}

func getCert(url string) (cert Cert, error error) {
	readData, writeData := io.Pipe()
	tokenString, err := util.ReturnToken()

	if err != nil {
		return cert, errors.New("Please login before attempting to obtain an SSH cert.")
	}

	var token = sshLogin{
		tokenString,
	}

	go func() {
		json.NewEncoder(writeData).Encode(token)
		writeData.Close()
	}()

	resp, err := http.Post(url+"/ssh", "application/json", readData)

	if err != nil {
		return cert, err
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&cert)

	return cert, nil
}

func init() {
	rootCmd.AddCommand(sshCmd)
	sshCmd.Flags().StringP("url", "u", "https://sshackd.herokuapp.com", "Set an alternate API URL")
	sshCmd.Flags().BoolP("print", "p", false, "Print cert and private key instead of adding to ssh-agent")
	sshCmd.Flags().BoolP("write", "w", false, "Write cert (key-cert.pub) and key (key) to files instead of adding them to ssh-agent")
}
