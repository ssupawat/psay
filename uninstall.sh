#!/bin/sh
# psay uninstaller: reverses install.sh (binary + kokoro venv + model).
# Usage: curl -fsSL https://raw.githubusercontent.com/ssupawat/psay/main/uninstall.sh | sh
# Keeps ~/.psay/settings.json (your lexicon). Pass --purge to remove it too.
set -eu

BIN_DIR=${PSAY_BIN_DIR:-"$HOME/.local/bin"}
purge=0
[ "${1:-}" = "--purge" ] && purge=1

if [ -f "$BIN_DIR/psay" ]; then
    rm -f "$BIN_DIR/psay"
    echo "removed $BIN_DIR/psay"
else
    echo "no psay binary at $BIN_DIR/psay"
fi

uv tool uninstall piper-tts >/dev/null 2>&1 || true
rm -rf "$HOME/.psay/venv"
echo "removed $HOME/.psay/venv (kokoro)"

if [ -d "$HOME/.psay/kokoro" ]; then
    rm -rf "$HOME/.psay/kokoro"
    echo "removed $HOME/.psay/kokoro (model)"
fi

if [ -d "$HOME/.psay/voices" ]; then
    rm -rf "$HOME/.psay/voices"
    echo "removed $HOME/.psay/voices"
fi

if [ "$purge" -eq 1 ]; then
    rm -rf "$HOME/.psay"
    echo "purged $HOME/.psay"
else
    echo "kept $HOME/.psay/settings.json (your lexicon) — delete $HOME/.psay or pass --purge to remove"
fi
