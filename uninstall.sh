#!/bin/sh
# psay uninstaller: reverses install.sh (binary + piper + voice).
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

if command -v uv >/dev/null 2>&1; then
    uv tool uninstall piper-tts >/dev/null 2>&1 || true
fi
if command -v pipx >/dev/null 2>&1; then
    pipx uninstall piper-tts >/dev/null 2>&1 || true
fi
echo "uninstalled piper-tts (if it was present)"

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
