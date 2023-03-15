#!/bin/bash
# diode agent binary location. by default, matches diode-agent container (see Dockerfile)
diode_agent_bin="${DIODE_AGENT_BIN:-/usr/local/bin/diode-agent}"
#
echo "Diode Agent Starting"
if [ $# -eq 0 ]; then
  exec "$diode_agent_bin" run
else
  exec "$diode_agent_bin" "$@"
fi
