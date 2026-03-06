# AIF (Agent Interfaces Framework) 🤖💻

[🇰🇷 한국어 버전 읽기 (Read in Korean)](README_ko.md)

> **Stop building heavy MCP servers. Give AI Agents what they really want: A fast, composable, single-binary CLI.**

## The Problem
Right now, the AI agent ecosystem is obsessed with MCP (Model Context Protocol). While MCP is great for some things, it requires running heavy background daemons, managing Node.js/Python dependencies, and dealing with complex network protocols.

When AI agents need to get things done on a local machine, they don't want to talk to a slow web server. They want to use the terminal. They want to pipe `|` outputs. They want **UNIX philosophy**.

## The Solution: AIF
**AIF (Agent Interfaces Framework)** is an open-source engine that auto-generates AI-friendly CLI tools from simple JSON/YAML API specifications. 

Instead of writing a custom MCP server for every SaaS tool, AIF advocates for building lightweight Go CLIs that output 100% pure JSON, making them instantly usable by agents like Claude Code, OpenClaw, and GitHub Copilot.

### The AI Agent Workflow
1. An AI Agent reads API documentation for a service.
2. The Agent creates a simple `spec.json` defining the endpoints and flags.
3. The Agent runs `aif build spec.json`.
4. A flawless, compiled Go binary CLI is instantly generated, complete with automatic token authentication (`auth login`) and JSON formatting.

## Proof of Concept
Included in this repository is the AIF Engine PoC and a sample `upost-spec.json` covering the [Upload-Post API](https://upload-post.com).

### Try it out:
```bash
# Build the AIF engine
go build -o aif main.go

# Generate the target CLI (e.g., upost) from the spec
./aif build upost-spec.json

# Test the generated CLI
./upost help
./upost auth login --token "your-api-key"
./upost text --title "Hello World" --platform x --platform linkedin
```

## Join the Movement
Initiated by a non-developer who just wanted AI to work better. Built with AI.
We are looking for Go and Rust developers who believe in building fast, stateless, and composable tools for the AI era to help expand this engine!
