/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/orb-community/diode/agent"
	"github.com/orb-community/diode/agent/config"
	"github.com/orb-community/diode/buildinfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	cfgFiles   []string
	Debug      bool
	OutputType string
	OutputPath string
	Host       string
	Port       uint32
)

func Version(cmd *cobra.Command, args []string) {
	fmt.Printf("diode-agent %s\n", buildinfo.GetVersion())
	os.Exit(0)
}

func Run(cmd *cobra.Command, args []string) {

	initConfig()

	// configuration
	var config config.Config

	config.Version = buildinfo.GetVersion()
	config.DiodeAgent.DiodeConfig.Debug = Debug
	config.DiodeAgent.DiodeConfig.Host = Host
	config.DiodeAgent.DiodeConfig.Port = strconv.FormatUint(uint64(Port), 10)
	config.DiodeAgent.DiodeConfig.OutputType = OutputType
	config.DiodeAgent.DiodeConfig.OutputPath = OutputPath

	err := viper.Unmarshal(&config)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("agent start up error (config): %w", err))
		os.Exit(1)
	}

	// logger
	var logger *zap.Logger
	atomicLevel := zap.NewAtomicLevel()
	if Debug {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else {
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
		_ = logger.Sync()
	}(logger)

	// new agent
	a, err := agent.New(logger, config)
	if err != nil {
		logger.Error("agent start up error", zap.Error(err))
		os.Exit(1)
	}

	// handle signals
	done := make(chan bool, 1)
	rootCtx, cancelFunc := context.WithCancel(context.WithValue(context.Background(), "routine", "mainRoutine"))

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		for {
			select {
			case <-sigs:
				logger.Warn("stop signal received stopping agent")
				a.Stop(rootCtx)
				cancelFunc()
			case <-rootCtx.Done():
				logger.Warn("mainRoutine context cancelled")
				done <- true
				return
			}
		}
	}()

	// start agent
	err = a.Start(rootCtx, cancelFunc)
	if err != nil {
		logger.Error("agent startup error", zap.Error(err))
		os.Exit(1)
	}

	<-done
}

func mergeOrError(path string) {

	v := viper.New()
	if len(path) > 0 {
		v.SetConfigFile(path)
		v.SetConfigType("yaml")
	}

	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	// note: viper seems to require a default (or a BindEnv) to be overridden by environment variables
	v.SetDefault("diode.config.debug", false)
	v.SetDefault("diode.config.output_type", "file")
	v.SetDefault("diode.config.output_path", "")
	v.SetDefault("diode.config.output_auth", "")

	if len(path) > 0 {
		cobra.CheckErr(v.ReadInConfig())
	}

	var fZero float64

	// check that version of config files are all matched up
	if versionNumber1 := viper.GetFloat64("version"); versionNumber1 != fZero {
		versionNumber2 := v.GetFloat64("version")
		if versionNumber2 == fZero {
			cobra.CheckErr("Failed to parse config version in: " + path)
		}
		if versionNumber2 != versionNumber1 {
			cobra.CheckErr("Config file version mismatch in: " + path)
		}
	}

	cobra.CheckErr(viper.MergeConfigMap(v.AllSettings()))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	mergeOrError("")
	if len(cfgFiles) == 0 {
		cobra.CheckErr("Config files were not provided")
	} else {
		for _, conf := range cfgFiles {
			mergeOrError(conf)
		}
	}
}

func main() {

	rootCmd := &cobra.Command{
		Use: "diode-agent",
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run diode-agent",
		Long:  `Run diode-agent`,
		Run:   Run,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show agent version",
		Run:   Version,
	}

	runCmd.Flags().StringSliceVarP(&cfgFiles, "config", "c", []string{}, "Path to config files (may be specified multiple times)")
	runCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Enable verbose (debug level) output")
	runCmd.PersistentFlags().StringVarP(&Host, "host", "i", "localhost", "Define agent server host")
	runCmd.PersistentFlags().Uint32VarP(&Port, "port", "p", 10911, "Define agent server port")
	runCmd.PersistentFlags().StringVarP(&OutputType, "type", "t", "file", "Define agent output type")
	runCmd.PersistentFlags().StringVarP(&OutputPath, "output_path", "o", "", "Define agent output path")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}
