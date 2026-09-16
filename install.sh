#!/bin/sh
# psay installer: binary + piper + voice in one shot (macOS).
# Usage: curl -fsSL https://raw.githubusercontent.com/ssupawat/psay/main/install.sh | sh
set -eu

REPO=ssupawat/psay
VOICE_DIR="$HOME/.psay/kokoro"
BIN_DIR=${PSAY_BIN_DIR:-"$HOME/.local/bin"}
KOKORO_BASE="https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files-v1.1"

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

echo "==> Setting up kokoro-onnx"
if command -v uv >/dev/null 2>&1; then
    uv venv "$HOME/.psay/venv" --python 3.12
    uv pip install --python "$HOME/.psay/venv/bin/python" kokoro-onnx soundfile
elif command -v python3 >/dev/null 2>&1; then
    python3 -m venv "$HOME/.psay/venv"
    "$HOME/.psay/venv/bin/pip" install --quiet kokoro-onnx soundfile
else
    echo "Need uv or python3 to set up kokoro" >&2
    exit 1
fi

echo "==> Downloading kokoro model (one-time, ~330MB)"
mkdir -p "$VOICE_DIR"
[ -f "$VOICE_DIR/kokoro-v1.0.onnx" ] || fetch "$KOKORO_BASE/kokoro-v1.0.onnx" "$VOICE_DIR/kokoro-v1.0.onnx"
[ -f "$VOICE_DIR/voices-v1.0.bin" ] || fetch "$KOKORO_BASE/voices-v1.0.bin" "$VOICE_DIR/voices-v1.0.bin"

case ":$PATH:" in
*":$BIN_DIR:"*) ;;
*) echo "NOTE: add $BIN_DIR to PATH" ;;
esac
echo "==> Done. Try: psay 'Hello from psay.'"
