#!/bin/bash

KNDP_CLI_URL="https://raw.githubusercontent.com/web-seven/kndp/release/0.1/scripts/kndp-cli.sh"
INSTALL_DIR="/usr/local/bin"

curl -sL "$KNDP_CLI_URL" -o "$INSTALL_DIR/kndp"
chmod +x "$INSTALL_DIR/kndp"

echo "KNDP CLI installed successfully. Run 'kndp --help' for usage information."