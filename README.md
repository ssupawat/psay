# psay

Speakable announcements for AI agents. psay turns developer jargon into short, speakable text and plays it through local TTS.

`psay "Found deadlock in auth_service.go due to A-B-A problem in the mutex pool."`

→ Kokoro speaks: *"Found deadlock in auth service. Go. Due to A B A problem in the mutex pool."*

Zero Go dependencies — stdlib only, plus a small Python venv (kokoro-onnx) and `afplay`.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/ssupawat/psay/main/install.sh | sh
```

Installs the `psay` binary to `~/.local/bin`, a kokoro-onnx venv at `~/.psay/venv`, and the Kokoro model (~330MB, one-time). Requires macOS. Or with Go: `go install github.com/ssupawat/psay@latest`. First announce: `psay 'Hello from psay.'` — config auto-initializes.

Uninstall with the matching script (keeps your lexicon unless `--purge`):

```bash
curl -fsSL https://raw.githubusercontent.com/ssupawat/psay/main/uninstall.sh | sh
```

## Problem

Agent announces are spoken for human ears. The listening test showed kokoro's espeak-lineage front-end already reads most developer text correctly — camelCase (`getUserById` → "get user by I-D"), symbols (`x != y` → "ex not-equals why"), extensions (`auth_service.go` → "auth service dot go"), versions (`v1.2.3` → "vee one point two point three").

What it gets wrong is a narrow class of word-level misreads:

- Lowercase acronyms read as words: `aws` → "awz"
- Dev shorthand: `db` → "dee bee"
- psay's own name: `psay` → "say"

All are word substitutions — exactly what a lexicon fixes. No phonetic engine needed. Brevity is not psay's problem either — the calling agent condenses its own announce before invoking psay (see [PSAY.md](./PSAY.md)).

## What psay does

psay sits between an AI agent and the TTS engine, running a text pipeline that transforms jargon before synthesis:

```mermaid
flowchart LR
  A["psay 'text'"] --> B["Lexicon replace
~/.psay/settings.json"]
  B --> E["kokoro
af_heart → wav"]
  E --> F["afplay"]
```

The listening test killed the planned identifier splitter and symbol mapper — the G2P front-end already does that work. What remains is the lexicon and the plumbing.

## Agent protocol

[PSAY.md](./PSAY.md) is the announcement protocol for AI agents. Run `psay init` once to install it at user level — it appends the protocol to the global instruction files that already exist: Codex (`~/.codex/AGENTS.md`), Claude (`~/.claude/CLAUDE.md`), pi (`~/.pi/agent/AGENTS.md`), ZCode (`~/.zcode/AGENTS.md`). It never creates files, so harnesses you don't use are skipped; idempotent, safe to re-run. The whole protocol in one line:

```bash
psay '[Action/Discovery] on [Target] because [Context]. [My Take / Advice]'
```

## Dependencies

Zero Go dependencies — everything is stdlib. External tools: kokoro-onnx (TTS, in a venv) and `afplay` (playback).

| Pipeline step | Implementation | Why |
|---|---|---|
| CLI parsing | stdlib `flag` package | `-v` voice override |
| Home dir path | stdlib `os.UserHomeDir()` | Available since Go 1.12 |
| TTS engine | [`thewh1teagle/kokoro-onnx`](https://github.com/thewh1teagle/kokoro-onnx) — `pip install kokoro-onnx soundfile` in `~/.psay/venv` | Local, offline neural TTS (Kokoro-82M). Model + voices downloaded once to `~/.psay/kokoro`; psay execs a 7-line embedded script |
| Playback | macOS `afplay` | Built-in, plays the wav |

## Configuration

Settings live at `~/.psay/settings.json`:

```json
{
  "voice": "af_heart",
  "speed": 1.0,
  "lexicon": {
    "psay": "p say",
    "aws": "A W S",
    "db": "database",
    "err": "error",
    "nil": "null",
    "repo": "repository",
    "auth": "authentication"
  }
}
```

The config auto-initializes if missing; a pre-kokoro file (with `length_scale`) migrates automatically. Voices: `af_heart` (best), `af_bella`, `am_michael`, and more in [Kokoro-82M/VOICES.md](https://huggingface.co/hexgrad/Kokoro-82M/blob/main/VOICES.md). `speed` is kokoro's pace — higher is faster. Lexicon is a simple key-value map — edit directly in a text editor. New misreads found in daily use become lexicon entries; code changes only if a whole token class proves broken.

