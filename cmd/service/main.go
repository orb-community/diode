/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/orb-community/diode/service"
	"github.com/orb-community/diode/service/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	envPrefix = "diode_service"
)

func main() {
	svcCfg := config.LoadConfig(envPrefix)

	// main logger
	var logger *zap.Logger
	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToLower(svcCfg.Base.LogLevel) {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		os.Stdout,
		atomicLevel,
	)
	logger = zap.New(core, zap.AddCaller())
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(logger) // flushes buffer, if any

	db := startSqliteDb(logger)

	svc := service.New(logger, svcCfg, db)
	defer func(svc service.Service) {
		err := svc.Stop()
		if err != nil {
			log.Fatalf("fatal error in stop the service: %e", err)
		}
	}(svc)

	errs := make(chan error, 2)

	err := svc.Start()
	if err != nil {
		logger.Error("unable to start agent data consumption", zap.Error(err))
		os.Exit(1)
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error("diode service terminated", zap.Error(err))
}

func startSqliteDb(logger *zap.Logger) (db *sql.DB) {
	db, err := sql.Open("sqlite3", ":memory")
	if err != nil {
		logger.Fatal("SQLite could not be initialized", zap.Error(err))
	}

	createInterfacesTableStatement, err := db.Prepare("CREATE TABLE IF NOT EXISTS interfaces (id TEXT PRIMARY KEY, policy TEXT, interfaceData TEXT )")
	if err != nil {
		logger.Fatal("error preparing interfaces statement", zap.Error(err))
	}
	_, err = createInterfacesTableStatement.Exec()
	if err != nil {
		logger.Fatal("error creating interfaces table", zap.Error(err))
	}
	logger.Debug("successfully created Interfaces table")
	createDeviceTableStatement, err := db.Prepare("CREATE TABLE IF NOT EXISTS devices (id TEXT PRIMARY KEY, policy TEXT, deviceData TEXT )")
	if err != nil {
		logger.Fatal("error preparing devices statement ", zap.Error(err))
	}
	_, err = createDeviceTableStatement.Exec()
	if err != nil {
		logger.Fatal("error creating devices table", zap.Error(err))
	}
	logger.Debug("successfully created devices table")

	return
}
