#!/bin/sh
# diode agent binary location. by default, matches diode-agent container (see Dockerfile)
diode_agent_bin="${DIODE_AGENT_BIN:-/usr/local/bin/diode-agent}"
#
if [ $# -eq 0 ]; then
  "$diode_agent_bin" run &
  echo $! > /var/run/diode-agent.pid
else
  "$diode_agent_bin" "$@" &
  echo $! > /var/run/diode-agent.pid
fi
