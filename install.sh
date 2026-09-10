#!/bin/sh
# psay installer: binary + piper + voice in one shot (macOS).
# Usage: curl -fsSL https://raw.githubusercontent.com/ssupawat/psay/main/install.sh | sh
set -eu

REPO=ssupawat/psay
VOICE=en_US-lessac-medium
BIN_DIR=${PSAY_BIN_DIR:-"$HOME/.local/bin"}

case "$(uname -s)" in
Darwin) ;;
*) echo "psay needs macOS (afplay). Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
arm64 | aarch64) arch=arm64 ;;
x86_64) arch=amd64 ;;
*) echo "Unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac

if command -v curl >/dev/null 2>&1; then
    fetch() { curl -fsSL -o "$2" "$1"; }
else
    fetch() { wget -qO "$2" "$1"; }
fi

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "==> Downloading psay (latest release)"
fetch "https://api.github.com/repos/$REPO/releases/latest" "$tmp/release.json"
tag=$(sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' "$tmp/release.json" | head -n1)
[ -n "$tag" ] || { echo "Could not resolve latest release" >&2; exit 1; }
fetch "https://github.com/$REPO/releases/download/$tag/psay_${tag#v}_darwin_${arch}.tar.gz" "$tmp/psay.tar.gz"
tar -xzf "$tmp/psay.tar.gz" -C "$tmp"
mkdir -p "$BIN_DIR"
mv "$tmp/psay" "$BIN_DIR/psay"
chmod +x "$BIN_DIR/psay"
echo "    installed $BIN_DIR/psay ($tag)"

echo "==> Installing piper"
if command -v uv >/dev/null 2>&1; then
    uv tool install --force piper-tts
elif command -v pipx >/dev/null 2>&1; then
    pipx install --force piper-tts
else
    echo "Need uv or pipx to install piper (brew install uv)" >&2
    exit 1
fi

echo "==> Downloading voice $VOICE"
mkdir -p "$HOME/.psay/voices"
if ! uvx --from piper-tts python -m piper.download_voices "$VOICE" --download-dir "$HOME/.psay/voices"; then
    base="https://huggingface.co/rhasspy/piper-voices/resolve/v1.0.0/en/en_US/lessac/medium"
    fetch "$base/$VOICE.onnx" "$HOME/.psay/voices/$VOICE.onnx"
    fetch "$base/$VOICE.onnx.json" "$HOME/.psay/voices/$VOICE.onnx.json"
fi

case ":$PATH:" in
*":$BIN_DIR:"*) ;;
*) echo "NOTE: add $BIN_DIR to PATH" ;;
esac
echo "==> Done. Try: psay 'Hello from psay.'"
