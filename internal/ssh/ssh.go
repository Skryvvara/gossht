package ssh

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/term"
)

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func SSHConnect(host, user string) {
	sshPath := path.Join(os.Getenv("HOME"), ".ssh")
	//TODO: support other keys than id_rsa and custom paths
	//keyPath := path.Join(sshPath, "id_rsa")

	var err error
	//var signer ssh.Signer

	// Read the private key file
	/*
		pKey, err := os.ReadFile(keyPath)
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
		}*/

	/*
		signer, err = ssh.ParsePrivateKey(pKey)
		if err != nil {
			fmt.Println(err.Error())
		}*/

	var hostkeyCallback ssh.HostKeyCallback
	hostkeyCallback, err = knownhosts.New(path.Join(sshPath, "known_hosts"))
	if err != nil {
		fmt.Println(err.Error())
	}

	conf := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: hostkeyCallback,
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
		Timeout: time.Millisecond * 1000,
	}

	client, err := ssh.Dial("tcp", host, conf)
	if err != nil {
		// Handle specific errors
		if isConnectionError(err) {
			fmt.Printf("Failed to dial SSH connection: %v\n", err)
			return
		}
		fmt.Printf("Unknown error while dialing SSH connection: %v\n", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Get the terminal file descriptor
	fd := int(os.Stdin.Fd())

	// Put the terminal into raw mode
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatalf("Failed to set terminal to raw mode: %v", err)
	}
	defer term.Restore(fd, oldState)

	// Get the terminal size
	width, height, err := term.GetSize(fd)
	if err != nil {
		log.Fatalf("Failed to get terminal size: %v", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.ECHOCTL:       0,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request a pseudo-terminal
	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		log.Fatalf("Request for pseudo terminal failed: %v", err)
	}

	// Set input and output
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	// Handle terminal resizing
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGWINCH)
		for range sigCh {
			width, height, _ := term.GetSize(fd)
			session.WindowChange(height, width)
		}
	}()

	// Handle termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh // Wait for signal
		session.Close()
	}()

	err = session.Wait()
	if err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
}

// Helper function to check if the error is due to connection issues
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific errors indicating connection issues
	if strings.Contains(err.Error(), "connection refused") {
		return true
	}

	// Check for network errors
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if se.Err == syscall.ECONNREFUSED || se.Err == syscall.EHOSTUNREACH || se.Err == syscall.ETIMEDOUT {
				return true
			}
		}
	}

	// Check for context deadline exceeded error
	if strings.Contains(err.Error(), "context deadline exceeded") {
		return true
	}

	return false
}
