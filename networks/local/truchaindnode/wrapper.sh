#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/truchaind/${BINARY:-truchaind}
ID=${ID:-0}
LOG=${LOG:-truchaind.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'truchaind' E.g.: -e BINARY=truchaind_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export TRUCHAINDHOME="/truchaind/node${ID}/truchaind"

if [ -d "`dirname ${TRUCHAINDHOME}/${LOG}`" ]; then
  "$BINARY" --home "$TRUCHAINDHOME" "$@" | tee "${TRUCHAINDHOME}/${LOG}"
else
  "$BINARY" --home "$TRUCHAINDHOME" "$@"
fi

chmod 777 -R /truchaind

