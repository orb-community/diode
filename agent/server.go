/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"context"
	"io"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/orb-community/diode/agent/backend"
	"github.com/orb-community/diode/agent/backend/factory"
	"github.com/orb-community/diode/agent/config"
	"gopkg.in/yaml.v3"
)

type ReturnValue struct {
	Message string `json:"message"`
}

func (a *diodeAgent) startServer(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)
	a.router = gin.New()

	a.router.Use(ginzap.Ginzap(a.logger, time.RFC3339, true))
	a.router.Use(ginzap.RecoveryWithZap(a.logger, true))

	a.router.GET("/api/v1/status", a.getStatus)
	a.router.GET("/api/v1/policies", a.getPolicies)
	a.router.POST("/api/v1/policies", a.createPolicy)
	a.router.GET("/api/v1/policies/:policy", a.getPolicy)
	a.router.DELETE("/api/v1/policies/:policy", a.deletePolicy)

	go func() {
		a.logger.Info("starting diode-agent server at: " + a.addr)
		if err := a.router.Run(a.addr); err != nil {
			a.Stop(ctx)
		}
	}()
	return nil
}

func (a *diodeAgent) getStatus(c *gin.Context) {
	a.stat.UpTime = time.Since(a.stat.StartTime)
	c.IndentedJSON(http.StatusOK, a.stat)
}

func (a *diodeAgent) getPolicies(c *gin.Context) {
	policies := make([]string, 0, len(a.policies))
	for k := range a.policies {
		policies = append(policies, k)
	}
	c.IndentedJSON(http.StatusOK, policies)
}

func (a *diodeAgent) getPolicy(c *gin.Context) {
	policy := c.Param("policy")
	rInfo, ok := a.policies[policy]
	if ok {
		c.YAML(http.StatusOK, rInfo.cf)
	} else {
		c.JSON(http.StatusNotFound, ReturnValue{"policy not found"})
	}
}

func (a *diodeAgent) createPolicy(c *gin.Context) {
	if t := c.Request.Header.Get("Content-type"); t != "application/x-yaml" {
		c.JSON(http.StatusForbidden, ReturnValue{"invalid Content-Type. Only 'application/x-yaml' is supported"})
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusForbidden, ReturnValue{err.Error()})
		return
	}
	var payload map[string]config.Policy
	if err = yaml.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusForbidden, ReturnValue{err.Error()})
		return
	}
	if len(payload) > 1 {
		c.JSON(http.StatusForbidden, ReturnValue{"only single policy allowed per request"})
		return
	}
	var policy string
	var data config.Policy
	for policy, data = range payload {
		_, ok := a.policies[policy]
		if ok {
			c.JSON(http.StatusConflict, ReturnValue{"policy already exists"})
			return
		}
		if len(data.Data) == 0 {
			c.JSON(http.StatusForbidden, ReturnValue{"data field is required"})
			return
		}
	}

	be, err := factory.GetBackend(data.Backend)
	if err != nil {
		c.JSON(http.StatusForbidden, ReturnValue{err.Error()})
		return
	}

	if data.Kind != Kind {
		c.JSON(http.StatusForbidden, ReturnValue{"invalid policy kind"})
		return
	}
	if err = be.Configure(a.logger, policy, a.pusher.GetChannel(), data.Data, data.Config); err != nil {
		c.JSON(http.StatusForbidden, ReturnValue{err.Error()})
		return
	}
	backendCtx := context.WithValue(a.ctx, "routine", policy)

	if err := be.Start(context.WithCancel(backendCtx)); err != nil {
		c.JSON(http.StatusForbidden, ReturnValue{err.Error()})
		return
	}
	a.policies[policy] = backendInfo{
		be: be,
		st: &backend.State{
			Status:        backend.Unknown,
			LastRestartTS: time.Now(),
		},
		cf: data,
	}
	c.YAML(http.StatusCreated, data)
}

func (a *diodeAgent) deletePolicy(c *gin.Context) {
	policy := c.Param("policy")
	r, ok := a.policies[policy]
	if ok {
		if err := r.be.Stop(a.ctx); err != nil {
			c.JSON(http.StatusForbidden, ReturnValue{err.Error()})
			return
		}
		delete(a.policies, policy)
		c.JSON(http.StatusOK, ReturnValue{policy + " was deleted"})
	} else {
		c.JSON(http.StatusNotFound, ReturnValue{"policy not found"})
	}
}
