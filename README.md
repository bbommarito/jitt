# jitt

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

It's early days, so there’s no fancy install script yet. For now:

```bash
go install github.com/bbommarito/jitt@latest
```

---

## 🧪 Testing

This project uses [Testify](https://github.com/stretchr/testify) for assertions and test helpers.

To run the tests:

```bash
go test ./...
```

---

## ❤️ Contributions welcome

Whether you're a seasoned Gopher or just curious, feel free to suggest changes, report bugs, or discuss ideas. All voices welcome — especially if you bring patience and a sense of humor.

---

## 🪪 License

MIT — do what you want, be cool about it.