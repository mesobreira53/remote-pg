package ssh

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

// Attributes will be private

// SSHClient defines the interface for an SSH client
type SSHClient interface {
	Connect() error
	RunCommand(cmd string) (string, error)
	Close()
}

// SSHClientConfig defines the configuration for an SSH client
type SSHClientConfig struct {
	Server   string
	Port     int
	Username string
	Password string
}

// NewSSHClient creates a new SSH client
func NewSSHClient(config SSHClientConfig) SSHClient {
	return &sshClient{config: config}
}

// sshClient is an implementation of the SSHClient interface
type sshClient struct {
	config SSHClientConfig
}

// Connect establishes an SSH connection
func (s *sshClient) Connect() error {
	// Define server details
	server := "192.168.1.201:22" // Change this to your server address
	username := "root"
	password := "!Password123."

	// Configure SSH client
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Not recommended for production
	}

	// Establish SSH connection
	client, err := ssh.Dial("tcp", server, config)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	// Create an SSH session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	// Capture standard output
	output, err := session.CombinedOutput("ls -l") // Change the command as needed
	if err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}

	// Print output
	fmt.Println(string(output))
	return nil
}

// RunCommand executes a command on the remote server
func (s *sshClient) RunCommand(cmd string) (string, error) {
	return "", nil
}

// Close closes the SSH connection
func (s *sshClient) Close() {
}
