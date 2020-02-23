package auth

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// PrivateKey Loads a private and public key from "path" and returns a SSH ClientConfig to authenticate with the server
func PrivateKey(username string, privateKey string, keyCallBack ssh.HostKeyCallback) (client ssh.ClientConfig, err error) {
	var signer ssh.Signer
	if signer, err = ssh.ParsePrivateKey([]byte(privateKey)); err != nil {
		return
	}

	client = ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: keyCallBack,
	}

	return
}

// PrivateKeyWithPassphrase Creates the configuration for a client that authenticates with a password protected private key
func PrivateKeyWithPassphrase(username, privateKey, passphrase string, keyCallBack ssh.HostKeyCallback) (client ssh.ClientConfig, err error) {
	var signer ssh.Signer
	if signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(passphrase)); err != nil {
		return
	}

	client = ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: keyCallBack,
	}

	return
}

// SshAgent Creates a configuration for a client that fetches public-private key from the SSH agent for authentication
func SshAgent(username string, keyCallBack ssh.HostKeyCallback) (client ssh.ClientConfig, err error) {
	var socket = os.Getenv("SSH_AUTH_SOCK")

	var conn net.Conn
	if conn, err = net.Dial("unix", socket); err != nil {
		return
	}

	var agentClient = agent.NewClient(conn)
	client = ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: keyCallBack,
	}
	return
}

// PasswordKey Creates a configuration for a client that authenticates using username and password
func PasswordKey(username string, password string, keyCallBack ssh.HostKeyCallback) (client ssh.ClientConfig, err error) {
	client = ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: keyCallBack,
	}
	return
}
