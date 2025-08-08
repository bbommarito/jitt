# jitt

[![CI](https://github.com/bbommarito/jitt/workflows/CI/badge.svg)](https://github.com/bbommarito/jitt/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bbommarito/jitt)](https://goreportcard.com/report/github.com/bbommarito/jitt)
[![codecov](https://codecov.io/gh/bbommarito/jitt/branch/main/graph/badge.svg)](https://codecov.io/gh/bbommarito/jitt)
[![Go Reference](https://pkg.go.dev/badge/github.com/bbommarito/jitt.svg)](https://pkg.go.dev/github.com/bbommarito/jitt)
[![GitHub release](https://img.shields.io/github/release/bbommarito/jitt.svg)](https://github.com/bbommarito/jitt/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight wrapper for Git that adds project-specific enhancements — starting with support for Jira-aware commit workflows.

**jitt** stands for **"Jira + Git + Tiny Tooling"** — the idea is to build small, targeted features that help integrate common workflows (like Jira ticket linking or commit hygiene) without introducing heavyweight dependencies or complex configuration.

---

## ✨ What it does

Right now, jitt:

- Passes any unknown commands straight through to Git
- Adds a basic `jitt jira init` command, which creates a `.jira` file in the current Git repo
- Avoids accidental `.jira` creation outside a Git repo
- Lays the foundation for smarter Jira integration (like enforcing ticket prefixes in commit messages)

---

## 🔧 Why build this?

I wanted something that feels like Git, behaves like Git, but can gently nudge teams (and myself) into better habits — like associating commits with real tickets or ensuring clean logs. All without enforcing complex Git hooks or rewriting history.

And also? I wanted an excuse to learn Go.

---

## 🧠 Be gentle — I’m new to Go!

This is one of my first projects in Go, and I’m learning as I build. So if you spot something odd or non-idiomatic, feel free to open an issue or PR — just maybe... do it with a little kindness?

I'm here to learn, improve, and have fun along the way.

---

## 🚧 Roadmap

Planned features (each added carefully and test-first):

- ✅ `jitt jira init` command
- ⏳ Enforce ticket key pattern in commits (e.g., `ABC-123: message`)
- ⏳ Add `jitt jira validate` for pre-commit hooks
- ⏳ Configurable Jira key prefixes
- ⏳ Optional project scaffolding (`.jira.json`, `.gitignore`, etc.)

---

## 📦 Installation

### Quick Install

```bash
go install github.com/bbommarito/jitt@latest
```

### Development

For development, clone the repo and use the provided Makefile:

```bash
git clone https://github.com/bbommarito/jitt.git
cd jitt

# Set up development environment
make dev-setup

# Build the binary
make build

# Run tests
make test

# See all available commands
make help
```

---

## 🧪 Testing

This project uses [Ginkgo](https://github.com/onsi/ginkgo) for BDD-style testing with [Gomega](https://github.com/onsi/gomega) for assertions. We've migrated from Testify to provide better test organization and readability.

To run the tests:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run with race detection
go test -race ./...
```

Or use Ginkgo directly for even prettier output:

```bash
# Install Ginkgo CLI (optional)
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Run with Ginkgo
ginkgo -r -v
```

---

## ❤️ Contributions welcome

Whether you're a seasoned Gopher or just curious, feel free to suggest changes, report bugs, or discuss ideas. All voices welcome — especially if you bring patience and a sense of humor.

---

## 🪪 License

MIT — do what you want, be cool about it.