package main

import (
	"context"
	"fmt"
	"os/exec"
	"remote-ssh/pkg/persistence"
	"time"
)

type PGOperator struct {
	storage persistence.Storage
}

type Operator interface {
	reconcile()
	check() error
}

func NewPGOperator(db_name string) PGOperator {
	return PGOperator{storage: persistence.NewSQLStorage(db_name)}
}

// reconcileLoop simulates a reconcile function running in a separate goroutine.
func (p PGOperator) reconcile(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Exiting...")
			return
		default:
			err := p.check()
			if err != nil {
				fmt.Println("Error Checking pg instances: ", err)
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func (p PGOperator) check() error {
	p.storage.Init_instance_db("pg_instance")
	instanceMap, err := p.storage.Read_all_pg_instance()
	if err != nil {
		return err
	}
	for pg_name, _ := range instanceMap {
		cmdline := fmt.Sprintf("docker ps | grep %s", pg_name)
		cmd := exec.Command("sh", "-c", cmdline)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Erro ao validar a instance: ", pg_name)
		}

	}
	return nil
}
