package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/ramasauskas/ispbet/db"
	"github.com/ramasauskas/ispbet/server"
	"github.com/rs/zerolog"
	"github.com/swithek/sessionup/memstore"
)

func main() {
	log := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	mainLog := log.With().Str("goroutine", "main").Logger()

	dbLog := log.With().Str("goroutine", "db").Logger()
	database, err := db.NewDB("test.sql", dbLog)
	if err != nil {
		mainLog.Fatal().Err(err).Msg("cannot create database")
		return
	}

	sessionStore := memstore.New(time.Minute)

	mainLog.Info().Msg("started session store")

	dbAdapter := &serverDBAdapter{
		db: database,
	}

	srvLog := log.With().Str("goroutine", "server").Logger()
	srv := server.NewServer(8080, sessionStore, &serverBetAdapter{}, &dummyEmail{
		log: log.With().Str("goroutine", "email").Logger(),
	}, dbAdapter, srvLog)

	doneCh := make(chan struct{}, 1)
	interCh := make(chan os.Signal, 1)

	signal.Notify(interCh, os.Interrupt)

	go func() {
		defer close(doneCh)

		if err = srv.Run(); err != nil {
			mainLog.Error().Err(err).Msg("cannot run server")
		}
	}()

	select {
	case <-doneCh:
	case <-interCh:
		mainLog.Info().Msg("received interrupt signal, closing application")
	}

	if err = srv.Close(); err != nil {
		mainLog.Error().Err(err).Msg("cannot close server")
	}

	sessionStore.StopCleanup()

	mainLog.Info().Msg("stopped session store")

	if err = database.Close(); err != nil {
		mainLog.Error().Err(err).Msg("cannot close database")
	}

	mainLog.Info().Msg("application gracefully closed")
}
