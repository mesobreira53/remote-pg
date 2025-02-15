package main

import (
	"log/slog"
	"net/http"
	"os"
	"remote-ssh/pkg/persistence"
	"remote-ssh/pkg/pgdocker"
	zfs "remote-ssh/pkg/zfs"
	"strconv"
)

type application struct {
	logger  *slog.Logger
	zfs     zfs.Storage
	pgdock  pgdocker.PGServer
	storage persistence.Storage
}

func newApplication(zpoolName string) *application {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	zfs := zfs.NewZFSStorage(zpoolName)
	pgdock := pgdocker.NewPGContainer()
	storage := persistence.NewSQLStorage("persistence.db")

	return &application{
		logger:  logger,
		zfs:     zfs,
		pgdock:  pgdock,
		storage: storage,
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

func (app *application) pg_destroy(w http.ResponseWriter, r *http.Request) {
	pg_name := r.URL.Query().Get("pg_name")
	err := app.pgdock.DestroyPGServer(pg_name)
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
	_, pg_port, err := app.pgdock.CreatePGServer(pg_name, pg_version, pg_datadir)
	if err != nil {
		app.serverError(w, r, err)
	} else {
		pg_port_int, err := strconv.Atoi(pg_port)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		app.storage.Save_pg_intance(pg_name, pg_port_int)
		w.WriteHeader(http.StatusOK)
	}
}
