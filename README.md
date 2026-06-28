# Goscaff

<img width="3168" height="1344" alt="Goscaff Banner" src="https://github.com/user-attachments/assets/df715bfd-09d6-4d2a-b19e-929b4489cc0f" />

<p align="center">
  <strong>Instant Go project scaffolding CLI.</strong><br>
  Generate production-ready Go backend projects in seconds with an interactive wizard or automation-friendly flags.
</p>

---

## ✨ Features

* ⚡ **Interactive wizard** — generate projects through a simple step-by-step CLI
* 🚀 **Multiple frameworks** — Gin and Fiber
* 🗄️ **Database selection** — PostgreSQL and MySQL
* 🏗️ **Architecture options** — Modular, Layered, and Full Setup
* 🔌 **Optional Dependency Injection** — None, Uber Dig, or Uber Fx
* 📦 **Go Modules ready** — automatically generates `go.mod`
* 🔧 **Git initialized** — project is ready to use immediately
* 🎨 **Clean CLI output** — colorful progress, configuration summary, and next steps
* 🤖 **Automation friendly** — fully configurable using command-line flags

---

## Installation

### Using Go

```bash
go install github.com/nbintang/goscaff@latest
```

Make sure your Go binary directory is available in your `PATH`.

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Verify installation:

```bash
goscaff --help
```

---

## Quick Start

Create a new project:

```bash
goscaff new myapp
```

The interactive wizard will guide you through:

```
? Select framework
❯ Gin
  Fiber

? Select database
❯ PostgreSQL
  MySQL

? Select architecture
❯ Modular
  Layered
  Full Setup

? Dependency Injection
❯ None
  Uber Dig
  Uber Fx

? Module path
❯ myapp

Configuration

Project Name : myapp
Module Path  : myapp
Framework    : Gin
Database     : PostgreSQL
Architecture : Modular
DI           : Uber Fx

? Continue?
```

Then goscaff will automatically:

* Create the project directory
* Generate project files
* Create `go.mod`
* Run `go mod tidy`
* Initialize Git
* Display the next steps

---

## Non-Interactive Mode

Perfect for automation or CI.

Generate a Gin project using PostgreSQL, Modular architecture, and Uber Fx:

```bash
goscaff new myapp \
  --framework gin \
  --db postgres \
  --architecture modular \
  --di uber-fx
```

Specify a custom module path:

```bash
goscaff new myapp \
  --module github.com/username/myapp
```

You can also specify a template directly:

```bash
goscaff new myapp \
  --template gin-postgres-uber-fx-modular
```

---

## Available Options

### Frameworks

| Value   | Description     |
| ------- | --------------- |
| `gin`   | Gin Framework   |
| `fiber` | Fiber Framework |

---

### Databases

| Value      | Description |
| ---------- | ----------- |
| `postgres` | PostgreSQL  |
| `mysql`    | MySQL       |

---

### Architectures

| Value        | Description                             |
| ------------ | --------------------------------------- |
| `modular`    | Feature-based modular architecture      |
| `layered`    | Traditional layered architecture        |
| `full-setup` | Production-ready full project structure |

---

### Dependency Injection

| Value      | Description     |
| ---------- | --------------- |
| `none`     | No DI container |
| `uber-dig` | Uber Dig        |
| `uber-fx`  | Uber Fx         |

---

## Command Reference

```bash
goscaff new [project-name]

Flags:

--module          Go module path
--framework       gin | fiber
--db              postgres | mysql
--architecture    modular | layered | full-setup
--di              none | uber-dig | uber-fx
--template        Template ID
```

See all available options:

```bash
goscaff new --help
```

---

## Example

```bash
goscaff new awesome-api \
  --module github.com/acme/awesome-api \
  --framework fiber \
  --db postgres \
  --architecture modular \
  --di uber-fx
```

---

## Philosophy

goscaff follows a few simple principles:

* **Simple by default** — generate projects without unnecessary complexity.
* **Interactive first** — sensible defaults through an intuitive wizard.
* **Flexible** — choose only the technologies you need.
* **Automation friendly** — every interactive option has a corresponding CLI flag.
* **Production-ready** — generated projects follow clean and maintainable structures.

---

## Contributing

Contributions are welcome.

1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Open a Pull Request.

---

## License

MIT License
