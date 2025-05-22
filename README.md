# Divar Recommender

**Divar Recommender** is an addon for the Divar application that enhances the chat experience by recommending similar and better posts directly within the post chat using a chatbot interface.

---

## 🚀 Features

* Recommends similar and high-quality posts in real-time chat
* Lightweight and easy to integrate
* Configurable and extendable via `config.yaml`
* Optional hot reloading support for fast development

---

## 🛠 Installation

Follow these steps to set up the project:

### 1. Install Go

Make sure Go is installed. You can download it from: [https://golang.org/dl/](https://golang.org/dl/)

### 2. Install Dependencies

Navigate to the project root and run:

```bash
go mod tidy
```

### 3. Create Configuration File

Copy the example config file:

```bash
cp config.yaml.example config.yaml
```

### 4. Initialize Your Configuration

Open `config.yaml` and fill in the required configuration values based on your environment.

---

## 🔄 (Optional) Enable Hot Reloading with Air

For better development experience, you can use [Air](https://github.com/air-verse/air) for hot reloading:

### Install Air

```bash
go install github.com/air-verse/air@latest
```

### Initialize Air Configuration

```bash
air init
```

Then, modify the `.air.toml` file to use the following command:

```toml
cmd = "go build -o ./tmp/main.exe ./app"
```

This ensures your application is built properly with every change.

---

## 📂 Project Structure

```
DivarRecommender/
├── app/                # Main application cmd
├── internal/           # Application source codes
├── config.yaml         # Your project config file
├── config.yaml.example # Example configuration
├── go.mod              # Go module file
├── .air.toml           # Air hot reload configuration
└── README.md
```

---

## 📄 License

[MIT](LICENSE)
