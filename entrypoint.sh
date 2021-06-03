#!/bin/sh
# This script is used to run main program after istio-sidecar is ready
# enovy doc: https://www.envoyproxy.io/docs/envoy/latest/operations/admin#get%E2%80%93ready

MAIN_PROGRAM="$@"

echo "start main program: ${MAIN_PROGRAM}"
# Run main program
exec ${MAIN_PROGRAM}