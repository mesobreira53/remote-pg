package pgdocker

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
)

// User to create a postgres docker container
type PGSContainer struct{}

type PGServer interface {
	CreatePGServer(pg_name string, pg_version string, pg_datadir string) error
	DestroyPGServer(pg_name string) error
}

func NewPGContainer() PGSContainer {
	return PGSContainer{}
}

// GenerateRandomString securely generates a random string of n bytes
// Used for Admin password
func generateRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:n], nil
}

func (pg PGSContainer) DestroyPGServer(pg_name string) error {
	fmt.Println("Destroying postgres docker container", pg_name, "...")
	cmd_line := fmt.Sprintf("docker stop %s", pg_name)
	fmt.Println("Command:", cmd_line)
	cmd := exec.Command(cmd_line)
	// Set output to OS stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Execute command
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd_line = fmt.Sprintf("docker rm -f %s", pg_name)
	fmt.Println("Command:", cmd_line)
	cmd = exec.Command(cmd_line)
	// Execute command
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (pg PGSContainer) CreatePGServer(pg_name string, pg_version string, pg_datadir string) error {
	pg_pass, err := generateRandomString(20)
	if err != nil {
		return err
	}
	cmd_line := fmt.Sprintf("docker run -d --name %s -e POSTGRES_PASSWORD=%s -d -p 15432:5432 --restart=always -v %s:/var/lib/postgresql/data  postgres:%s", pg_name, pg_pass, pg_datadir, pg_version)
	cmd := exec.Command(cmd_line)
	// Set output to OS stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Execute command
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
