#!/bin/bash
set -e

APPCHAIN_NAME=$1
PLUGIN_CONFIG=$2

sidecar --repo=/root/sidecar appchain method register --admin-key ./key.json --method fabappchain \
     --doc-addr /ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi \
     --doc-hash QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi \
     --name "${APPCHAIN_NAME}" --type fabric --desc="test for fabric" --version v1.4.3 \
     --validators ./"${PLUGIN_CONFIG}"/fabric.validators --consensus raft


command1=$(sidecar --repo=/root/sidecar rule deploy --path ./"${PLUGIN_CONFIG}"/validating.wasm --method fabappchain --admin-key ./key.json)
address=$(echo "$command1"|grep -o '0x.\{40\}')
echo "${address}"

command2=$(sidecar --repo=/root/sidecar rule bind --addr "${address}" --method fabappchain --admin-key ./key.json)
proposalID=$(echo "$command2"|grep -o '0x.\{42\}')
echo "${proposalID}"

sidecar --repo=/root/sidecar start

exec "$@"