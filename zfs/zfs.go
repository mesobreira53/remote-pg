package zfs

import (
	"fmt"
	"os"
	"os/exec"
)

// Attributes will be private
type ZFSStorage struct {
	zfsPoolName string
}

type Storage interface {
	CreateZFS(volumeName string, size string) error
	CreateZFSSnapshot() error
}

func NewZFSStorage(poolName string) ZFSStorage {
	return ZFSStorage{zfsPoolName: poolName}
}

func (zs ZFSStorage) CreateZFS(volumeName string, size string) error {
	fmt.Println("Creating zfs volume", volumeName, "...")
	cmd_line := fmt.Sprintf("zfs create -o mountpoint=/mnt/%s, %s/%s", volumeName, zs.zfsPoolName, volumeName)
	fmt.Println("Command:", cmd_line)
	cmd := exec.Command(cmd_line)
	// Set output to OS stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Execute command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	} else {
		fmt.Println("ZFS volume created successfully!")
	}
	return nil
}

func (zs ZFSStorage) CreateZFSSnapshot() error {
	return nil
}
