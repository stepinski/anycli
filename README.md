# anycli

> Terminal client for AnythingLLM — zero-setup codebase chat, Unix-native piping, beautiful TUI.

[![CI](https://github.com/yourusername/anycli/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/anycli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/anycli)](https://goreportcard.com/report/github.com/yourusername/anycli)
[![Release](https://img.shields.io/github/v/release/yourusername/anycli)](https://github.com/yourusername/anycli/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Why anycli?

The [official AnythingLLM CLI](https://github.com/Mintplex-Labs/anything-llm-cli) requires Node.js/Bun. `anycli` is a single compiled binary — no runtime, no `npm install`, no version conflicts.

```bash
# official CLI
bun run start prompt "hello"

# anycli
anycli "hello"
```

|  | anycli | anything-llm-cli |
|--|--------|-----------------|
| Single binary | ✅ | ❌ (requires Bun) |
| Install via brew | ✅ | ❌ |
| Streaming | ✅ | ✅ |
| Pipe-first UX | ✅ | partial |
| Auto workspace from git | ✅ | ❌ |
| Interactive TUI | ✅ | ❌ |
| JSON output | ✅ | ❌ |

---

## Install

```bash
# Homebrew (macOS/Linux)
brew install yourusername/tap/anycli

# Go
go install github.com/yourusername/anycli@latest

# Download binary
# https://github.com/yourusername/anycli/releases
```

---

## Quickstart

```bash
# 1. Run AnythingLLM desktop — anythingllm.com/download
# 2. Get API key: Settings > Developer API
# 3. Configure anycli
anycli init

# 4. Chat
anycli "how does authentication work?"
```

---

## Usage

### Basic chat

```bash
# Ask anything — uses your default workspace
anycli "summarise my notes from this week"

# Use a specific workspace
anycli -w transcripts "what was decided in last week's call?"

# RAG-only mode (no chat history, just document retrieval)
anycli -m query "find everything about the landlord situation"
```

### Pipe anything in

```bash
# Code review
git diff | anycli "write a commit message"
git diff | anycli "review this for bugs"

# Log analysis
cat error.log | anycli "what's wrong here?"
kubectl logs mypod | anycli "is this normal?"

# Clipboard
pbpaste | anycli "summarise this"    # macOS
xclip -o | anycli "summarise this"  # Linux

# Capture output
anycli "write a README for this project" > README.md
```

### Auto workspace detection

```bash
# anycli detects your git repo name and uses it as workspace
cd ~/code/myproject
anycli "explain this codebase"
# → uses workspace "myproject" automatically
```

### Interactive TUI

```bash
anycli tui
```

### Upload documents

```bash
anycli upload meeting-notes.md
anycli upload --workspace transcripts call-2025-01-15.md
```

### Workspaces

```bash
anycli workspaces          # list all
anycli workspaces create   # create new
```

### Scripting

```bash
# JSON output for scripting
anycli --json "summarise" | jq '.response'

# No streaming, capture full response
SUMMARY=$(anycli -S "summarise this week")

# Exit codes: 0 = success, 1 = error, 2 = no results
if anycli -m query "critical bugs" > /dev/null; then
  echo "found relevant docs"
fi
```

---

## Configuration

```bash
anycli init   # interactive setup
```

Config lives at `~/.config/anycli/config.yaml`:

```yaml
url: http://localhost:3001
api_key: your-key-here
workspace: vault           # default workspace
stream: true
mode: chat

# Used by `anycli brief`
priorities:
  - "health (sleep >= 7h)"
  - "deep work"
  - "obligations"
```

All settings can be overridden with environment variables:

```bash
ANYCLI_WORKSPACE=transcripts anycli "last week's calls"
ANYCLI_URL=http://other-host:3001 anycli "question"
```

---

## `anycli brief`

Generates a daily brief using your priorities as a schema:

```bash
anycli brief
```

```
📋 Daily Brief — Tuesday 2025-03-04

STATUS
  ✅ Health — sleep 7.5h, exercise logged yesterday
  🔴 Deep work — no focused blocks this week
  🟡 Obligations — 2 items pending response

BLOCKERS
  Waiting on landlord response before next step

TODAY
  1. Block 2h deep work before noon
  2. Follow up on landlord email
  3. Finish auth PR
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). PRs welcome.

---

## Related

- [AnythingLLM](https://github.com/Mintplex-Labs/anything-llm) — the server this CLI talks to
- [anything-llm-cli](https://github.com/Mintplex-Labs/anything-llm-cli) — official Node.js CLI
