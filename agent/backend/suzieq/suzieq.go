/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package suzieq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/orb-community/diode/agent/backend"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var Tables = [...]string{"device", "interfaces", "inventory", "vlan"}

const PollerTable = "sqPoller"

type suzieqBackend struct {
	logger        *zap.Logger
	policyName    string
	returnPrefix  string
	inventoryPath string
	proc          *cmd.Cmd
	statusChan    <-chan cmd.Status
	pusher        chan []byte
	startTime     time.Time
	cancelFunc    context.CancelFunc
	ctx           context.Context
	channel       chan string
}

var _ backend.Backend = (*suzieqBackend)(nil)

func New() backend.Backend {
	return &suzieqBackend{channel: make(chan string, 8)}
}

func (s *suzieqBackend) getProcRunningStatus() (backend.RunningStatus, string, error) {
	status := s.proc.Status()
	if status.Error != nil {
		errMsg := fmt.Sprintf("suzieq process error: %v", status.Error)
		return backend.BackendError, errMsg, status.Error
	}
	if status.Complete {
		err := s.proc.Stop()
		return backend.Offline, "suzieq process ended", err
	}
	if status.StopTs > 0 {
		return backend.Offline, "suzieq process ended", nil
	}
	return backend.Running, "", nil
}

func (s *suzieqBackend) Configure(logger *zap.Logger, name string, pusher chan []byte, data map[string]interface{}) error {
	var prs bool
	var inventory interface{}
	if inventory, prs = data["inventory"]; !prs {
		return errors.New("you must set suzieq inventory")
	}

	d, err := yaml.Marshal(&inventory)
	if err != nil {
		return err
	}

	s.inventoryPath = "/tmp/" + name + "_inventory.yml"
	err = os.WriteFile(s.inventoryPath, d, 0644)
	if err != nil {
		return err
	}

	s.logger = logger
	s.policyName = name
	s.returnPrefix = "{\"" + name + "\":{\"backend\":\"suzieq\", "
	s.pusher = pusher

	return nil
}

func (s *suzieqBackend) Version() (string, error) {

	envCmd := cmd.NewCmd("sq-poller -v")
	status := <-envCmd.Start()
	if len(status.Stdout) == 0 {
		return "", errors.New("sq-poller not found")
	}
	return status.Stdout[0], nil
}

func (s *suzieqBackend) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	s.startTime = time.Now()
	s.cancelFunc = cancelFunc
	s.ctx = ctx

	sOptions := []string{
		"-I",
		s.inventoryPath,
		"-o",
		"logging",
		"--no-coalescer",
		"--run-once",
		"update",
	}

	s.logger.Info("suzieq startup", zap.Strings("arguments", sOptions))

	s.proc = cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, "sq-poller", sOptions...)
	s.statusChan = s.proc.Start()

	matchOutput := regexp.MustCompile(`\bsuzieq.poller.worker.writers.logging\b`)

	// log STDOUT and STDERR lines streaming from Cmd
	go func() {
		for s.proc.Stdout != nil || s.proc.Stderr != nil || s.proc.Done() != nil {
			select {
			case line, open := <-s.proc.Stdout:
				if !open {
					s.proc.Stdout = nil
					continue
				}
				s.channel <- line
			case line, open := <-s.proc.Stderr:
				if !open {
					s.proc.Stderr = nil
					continue
				}
				s.logger.Info("suzieq stderr", zap.String("log", line))
			case <-s.proc.Done():
				s.Stop(ctx)
				return
			}
		}
	}()

	//proccess Suzieq stdout concurrently
	go func() {
		for {
			select {
			case line := <-s.channel:
				if matchOutput.MatchString(line) {
					_, output, _ := strings.Cut(line, "{")
					s.proccessDiscovery(output)
				} else {
					s.logger.Info("suzieq stdout", zap.String("log", line))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// wait for simple startup errors
	time.Sleep(time.Second)

	status := s.proc.Status()

	if status.Error != nil {
		s.logger.Error("suzieq startup error", zap.Error(status.Error))
		return status.Error
	}

	if status.Complete {
		err := s.proc.Stop()
		if err != nil {
			s.logger.Error("proc.Stop error", zap.Error(err))
		}
		return errors.New("suzieq startup error, check log")
	}

	s.logger.Info("suzieq process started", zap.Int("pid", status.PID))

	return nil
}

func (s *suzieqBackend) proccessDiscovery(data string) {
	discoveryData := []byte(s.returnPrefix + data + "}")
	var jsonData map[string]map[string]interface{}
	if err := json.Unmarshal(discoveryData, &jsonData); err != nil {
		s.logger.Error("process suzieq output error", zap.Error(err))
		return
	}

	for k, v := range jsonData[s.policyName] {
		for _, d := range Tables {
			if k == d {
				s.pusher <- discoveryData
			}
		}
		if k == PollerTable {
			status := ""
			service := ""
			pollerData := v.([]interface{})
			for _, pollerField := range pollerData {
				for k, i := range pollerField.(map[string]interface{}) {
					if k == "service" {
						for _, d := range Tables {
							if i.(string) == d {
								service = d
								break
							}
						}
					} else if k == "status" {
						status = fmt.Sprintf("%v", i)
					}
				}
			}
			if service != "" && status != "0" {
				s.logger.Error("suzieq poller for service '" + service + "' failed with status: " + status)
			}
		}
	}
}

func (s *suzieqBackend) Stop(ctx context.Context) error {
	s.logger.Info("routine call to stop suzieq", zap.Any("routine", ctx.Value("routine")))
	defer s.cancelFunc()
	err := s.proc.Stop()
	finalStatus := <-s.statusChan
	if err != nil {
		s.logger.Error("suzieq shutdown error", zap.Error(err))
		return err
	}
	s.logger.Info("suzieq process stopped", zap.Int("pid", finalStatus.PID), zap.Int("exit_code", finalStatus.Exit))
	return nil
}

func (s *suzieqBackend) FullReset(ctx context.Context) error {
	return nil
}

func (s *suzieqBackend) GetStartTime() time.Time {
	return s.startTime
}

func (s *suzieqBackend) GetCapabilities() (map[string]interface{}, error) {
	//TODO: implement capabilities which probably will be
	jsonBody := make(map[string]interface{})
	return jsonBody, nil
}

func (s *suzieqBackend) GetRunningStatus() (backend.RunningStatus, string, error) {
	runningStatus, errMsg, err := s.getProcRunningStatus()
	if runningStatus != backend.Running {
		return runningStatus, errMsg, err
	}
	return runningStatus, "", nil
}
