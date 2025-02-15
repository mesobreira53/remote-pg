package main

import (
	"log/slog"
	"net/http"
	"os"
	"remote-ssh/pgdocker"
	zfs "remote-ssh/zfs"
)

type application struct {
	logger *slog.Logger
	zfs    zfs.Storage
	pgdock pgdocker.PGServer
}

func newApplication(zpoolName string) *application {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	zfs := zfs.NewZFSStorage(zpoolName)
	pgdock := pgdocker.NewPGContainer()

	return &application{
		logger: logger,
		zfs:    zfs,
		pgdock: pgdock,
	}
}

func (app *application) app_home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Handler function for the zfs_create route
func (app *application) zfs_vol_create(w http.ResponseWriter, r *http.Request) {
	volumeName := r.URL.Query().Get("volumeName")
	size := r.URL.Query().Get("size")
	err := app.zfs.CreateZFS(volumeName, size)
	if err != nil {
		app.serverError(w, r, err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (app *application) pg_create(w http.ResponseWriter, r *http.Request) {
	pg_name := r.URL.Query().Get("pg_name")
	pg_version := r.URL.Query().Get("pg_version")
	pg_datadir := r.URL.Query().Get("pg_datadir")
	err := app.pgdock.CreatePGServer(pg_name, pg_version, pg_datadir)
	if err != nil {
		app.serverError(w, r, err)
	} else {
		app.logger.Info("Postgres container created successfully")
		w.WriteHeader(http.StatusOK)
	}
}
