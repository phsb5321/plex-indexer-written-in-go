# 🎬 Plex Indexer Written in Go 📚

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/plex-indexer-written-in-go)](https://goreportcard.com/report/github.com/yourusername/plex-indexer-written-in-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📖 Table of Contents

- [🎬 Plex Indexer Written in Go 📚](#-plex-indexer-written-in-go-)
  - [📖 Table of Contents](#-table-of-contents)
  - [🚀 About](#-about)
  - [✨ Features](#-features)
  - [📥 Installation](#-installation)
  - [🛠️ Usage](#️-usage)
    - [1. Generate JSON 📊](#1-generate-json-)
    - [2. Create Symlinks 🔗](#2-create-symlinks-)
    - [3. Full Process (Generate \& Symlink) 🔄](#3-full-process-generate--symlink-)
    - [Flags 🚩](#flags-)
  - [⚙️ Configuration](#️-configuration)
  - [🤝 Contributing](#-contributing)
  - [📄 License](#-license)

## 🚀 About

Plex Indexer Written in Go is a powerful tool designed to organize and index your media files for optimal Plex library management. This Go-based application streamlines the process of structuring your media collection, making it easier to maintain and navigate within Plex.

## ✨ Features

- 🗂️ Generates a structured JSON representation of your media files
- 🔗 Creates symlinks for easy Plex library organization
- 🧹 Cleans and reorders season and episode naming
- 🎭 Supports various video file formats
- 🖥️ Command-line interface for easy integration into your workflow

## 📥 Installation

1. Ensure you have Go 1.16 or later installed on your system.

2. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/plex-indexer-written-in-go.git
   cd plex-indexer-written-in-go
   ```

3. Build the application:

   ```bash
   go build -o plex-indexer cmd/course-organizer/main.go
   ```

## 🛠️ Usage

The Plex Indexer offers three main commands:

### 1. Generate JSON 📊

```bash
./plex-indexer generate -c /path/to/media/directory -j output.json
```

### 2. Create Symlinks 🔗

```bash
./plex-indexer symlink -j input.json -o /path/to/output/directory
```

### 3. Full Process (Generate & Symlink) 🔄

```bash
./plex-indexer full -c /path/to/media/directory -o /path/to/output/directory -g
```

### Flags 🚩

- `-c, --course`: Specifies the media directory path
- `-o, --output`: Specifies the output directory path for symlinks
- `-j, --json`: Specifies the JSON file to read from or write to
- `-g, --generate-json`: (For `full` command only) Generates JSON files

## ⚙️ Configuration

No additional configuration is required. The application uses command-line flags for all necessary inputs.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤️ by [Pedro Balbino](https://github.com/phsb5321)
