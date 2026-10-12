package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func newCA(t *testing.T) ssh.Signer {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		t.Fatal(err)
	}

	return signer
}

func newCert(t *testing.T, ca ssh.Signer, id string) (*ssh.Certificate, *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	pub, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	cert := &ssh.Certificate{
		Key:             pub,
		KeyId:           id,
		CertType:        ssh.UserCert,
		ValidPrincipals: []string{"deploy"},
		ValidAfter:      uint64(time.Now().Add(-time.Minute).Unix()),
		ValidBefore:     uint64(time.Now().Add(time.Hour).Unix()),
	}

	if err := cert.SignCert(rand.Reader, ca); err != nil {
		t.Fatal(err)
	}

	return cert, key
}

// Talks to the keyring over the agent wire protocol, like sb does with ssh-agent.
func newAgentClient(t *testing.T) agent.ExtendedAgent {
	client, server := net.Pipe()
	keyring := agent.NewKeyring()

	go agent.ServeAgent(keyring, server)

	t.Cleanup(func() {
		client.Close()
		server.Close()
	})

	return agent.NewClient(client)
}

func certIDs(t *testing.T, sshAgent agent.Agent) []string {
	keys, err := sshAgent.List()
	if err != nil {
		t.Fatal(err)
	}

	var ids []string

	for _, k := range keys {
		pub, err := ssh.ParsePublicKey(k.Blob)
		if err != nil {
			t.Fatal(err)
		}

		if cert, ok := pub.(*ssh.Certificate); ok {
			ids = append(ids, cert.KeyId)
		} else {
			ids = append(ids, "plain-key")
		}
	}

	return ids
}

func TestReplaceCertRemovesOlderCertFromSameCA(t *testing.T) {
	sshAgent := newAgentClient(t)
	ca := newCA(t)

	first, firstKey := newCert(t, ca, "first")
	if err := replaceCert(sshAgent, first, firstKey, 3600); err != nil {
		t.Fatal(err)
	}

	second, secondKey := newCert(t, ca, "second")
	if err := replaceCert(sshAgent, second, secondKey, 3600); err != nil {
		t.Fatal(err)
	}

	ids := certIDs(t, sshAgent)
	if len(ids) != 1 || ids[0] != "second" {
		t.Fatalf("expected only the new cert in the agent, got %v", ids)
	}
}

func TestReplaceCertKeepsUnrelatedIdentities(t *testing.T) {
	sshAgent := newAgentClient(t)
	ca := newCA(t)

	plainKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	if err := sshAgent.Add(agent.AddedKey{PrivateKey: plainKey}); err != nil {
		t.Fatal(err)
	}

	other, otherKey := newCert(t, newCA(t), "other-ca")
	if err := sshAgent.Add(agent.AddedKey{PrivateKey: otherKey, Certificate: other}); err != nil {
		t.Fatal(err)
	}

	old, oldKey := newCert(t, ca, "old")
	if err := replaceCert(sshAgent, old, oldKey, 3600); err != nil {
		t.Fatal(err)
	}

	fresh, freshKey := newCert(t, ca, "fresh")
	if err := replaceCert(sshAgent, fresh, freshKey, 3600); err != nil {
		t.Fatal(err)
	}

	got := map[string]bool{}
	for _, id := range certIDs(t, sshAgent) {
		got[id] = true
	}

	if len(got) != 3 || !got["plain-key"] || !got["other-ca"] || !got["fresh"] {
		t.Fatalf("expected plain-key, other-ca and fresh to remain, got %v", got)
	}
}

func TestReplaceCertIntoEmptyAgent(t *testing.T) {
	sshAgent := newAgentClient(t)

	cert, key := newCert(t, newCA(t), "only")
	if err := replaceCert(sshAgent, cert, key, 3600); err != nil {
		t.Fatal(err)
	}

	ids := certIDs(t, sshAgent)
	if len(ids) != 1 || ids[0] != "only" {
		t.Fatalf("expected the new cert in the agent, got %v", ids)
	}
}
