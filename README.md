# goscaff

**Instant Go project scaffolding CLI.**

`goscaff` is a lightweight Go CLI tool that scaffolds clean, production-ready Go backend projects in seconds—without boilerplate fatigue or framework lock-in.

---

## ✨ Features

* ⚡ **Instant setup** — create a new Go project with one command
* 🧱 **Presets** — `base` (minimal) and `full` (production-ready)
* 🧩 **Flexible architecture** — no forced DI or framework coupling
* 📦 **Go modules ready** — `go.mod` generated automatically
* 🧰 **Git initialized** — repository ready out of the box
* 🌱 **Environment files included** — `.env`, `.env.example`, `.env.local`
* 🎨 **Clean CLI output** — readable, colored progress & next-steps

---

## 📦 Installation

### Using Go

```bash
go install github.com/nbintang/goscaff@latest
```

Ensure `$GOPATH/bin` or `$HOME/go/bin` is in your `PATH`.

---

## 🚀 Usage

### Create a new project (default: base preset)

```bash
goscaff new myapp
```

This will:

* Create a `myapp` directory
* Scaffold the **base** project structure
* Generate `go.mod`
* Run `go mod tidy`
* Initialize a git repository
* Print clear **next steps** to run the project

---

### Specify module path (optional)

```bash
goscaff new myapp --module github.com/username/myapp
```

If `--module` is omitted, the module name defaults to the project name.

---

### Use full preset

```bash
goscaff new myapp --preset full
```

The `full` preset includes additional infrastructure and production-oriented defaults.

---

### Choose database (full preset only)

```bash
goscaff new myapp --preset full --db mysql
```

Supported databases:

* `postgres` (default)
* `mysql`

> Database overlays are applied **only** for the `full` preset. The `base` preset stays minimal.

---

## 📂 Example Project Structure

```text
myapp/
├── cmd/
│   └── api/
│   │   └── main.go
│   └── migrate
│   │   └── main.go
│   └── ...
│       
├── internal/
│   ├── auth/
│   ├── apperr/
│   ├── infra/
│   │   ├── database/
│   │   ├── cache/
│   │   └── ...
│   ├── http/
│   └── ...
│
├── pkg/
│   ├── env/
│   ├── slice/
│   └── ...
├── .env
├── .env.example
├── .env.local
├── go.mod
└── README.md
```

---

## 🧭 Philosophy

`goscaff` is designed with a few simple principles:

* **Minimal by default** — start clean, add complexity only when needed
* **Fast feedback** — scaffolding should take seconds, not minutes
* **Structure without lock-in** — you own the architecture decisions

---

## 🛣️ Roadmap

* [ ] Interactive prompts (`goscaff new`)
* [ ] Preset selection wizard
* [ ] Custom template support
* [ ] Prebuilt binaries (Windows / macOS / Linux)

---

## 🤝 Contributing

Contributions are welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

---

## 📄 License

MIT License

---

## ⭐ Acknowledgements

* [Cobra](https://github.com/spf13/cobra) — CLI framework
* [fatih/color](https://github.com/fatih/color) — terminal coloring
* Go community for inspiring great too
