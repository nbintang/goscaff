# goscaff

**Instant Go project scaffolding.**

`goscaff` is an instant Go scaffolding CLI that helps you generate production-ready Go backend projects in seconds — without boilerplate fatigue.

---

## ✨ Features

* ⚡ **Instant project setup** — one command to get started
* 🧱 **Scaffolding presets** (`base`, `full`)
* 🔁 **Flexible architecture** — no forced DI or framework lock-in
* 📦 **Go modules ready** (`go.mod` generated)
* 🧰 **Git initialized** automatically
* 🧪 **Production-oriented structure** (`cmd/`, `internal/`, `pkg/`)

---

## 📦 Installation

### Using Go (recommended)

```bash
go install github.com/nbintang/goscaff@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`.

---

## 🚀 Usage

### Create a new project

```bash
goscaff new myapp
```

This will:

* Create a `myapp` directory
* Scaffold a clean Go project structure
* Generate `go.mod`
* Run `go mod tidy`
* Initialize a git repository

---

### Specify module path (optional)

```bash
goscaff new myapp --module github.com/username/myapp
```

If `--module` is not provided, the module name defaults to the project name.

---

### Choose database

```bash
goscaff new myapp --db mysql
```

Supported databases:

* `postgres` (default)
* `mysql`

---

## 📂 Project Structure

```
myapp/
├── cmd/
│   ├── api/
│   ├── migrate/
│   └── seed/
├── internal/
│   ├── user/
│   ├── auth/
│   ├── infra/
│   └── http/
├── pkg/
│   └── utils/
├── scripts/
│   ├── migrate.sh
│   └── seed.sh
├── go.mod
├── go.sum
└── README.md
```

---

## 🧭 Philosophy

`goscaff` is built with these principles in mind:

* **Instant, not complex** — reduce setup time, not add layers
* **Practical over opinionated** — structure is provided, decisions stay with you
* **Scalable by default** — simple to start, easy to extend

---

## 🛣️ Roadmap

* [ ] Interactive prompts (`goscaff new`)
* [ ] Preset selection (`base`, `full`)
* [ ] Custom template support
* [ ] Binary releases (Windows / macOS / Linux)

---

## 🤝 Contributing

Contributions are welcome!

1. Fork the repository
2. Create a new branch
3. Make your changes
4. Submit a pull request

---

## 📄 License

MIT License

---

## ⭐ Acknowledgements

* [Cobra](https://github.com/spf13/cobra) — CLI framework
* Go community for inspiring great tooling
