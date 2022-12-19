package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/ramasauskas/ispbet/autobet"
	"github.com/ramasauskas/ispbet/autoreport"
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

	betDBAdapter := &betDBAdapter{
		db: database,
	}

	better := &better{
		db: betDBAdapter,
	}

	betSrv := serverBetAdapter{
		better: better,
	}

	reportDB := &reportDB{
		db: database,
	}

	dummyEm := &dummyEmail{
		log: log.With().Str("goroutine", "auto_report").Logger(),
	}

	autoBetDB := &autobetDB{
		db:  database,
		bet: better,
	}

	autoOddsDB := &autoOddsDB{
		db: database,
	}

	autoOddsWorker := newOddsWorker(autoOddsDB, log.With().Str("goroutine", "auto_odds").Logger())

	go autoOddsWorker.work()

	mainLog.Info().Msg("started auto odds worker")

	autoWorker := autobet.NewWorker(autoBetDB, autoBetDB, log.With().Str("goroutine", "autobet").Logger())

	go autoWorker.Work()

	mainLog.Info().Msg("started auto better")

	reportWorker := autoreport.NewWorker(reportDB, dummyEm, log.With().Str("goroutine", "email").Logger())
	if err := reportWorker.Work(); err != nil {
		mainLog.Fatal().Err(err).Msg("cannot run worker")
	}

	mainLog.Info().Msg("started report worked")

	srvLog := log.With().Str("goroutine", "server").Logger()
	srv := server.NewServer(8080, sessionStore, &betSrv, &betSrv, dummyEm, dbAdapter, srvLog)

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

	reportWorker.Close()

	mainLog.Info().Msg("stopped report worker")

	autoOddsWorker.stop()

	mainLog.Info().Msg("stopped auto odds worker")

	mainLog.Info().Msg("application gracefully closed")
}
