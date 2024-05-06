# Quotify

Quotify is an app written in Go for managing quotes. It utilizes a JSON file-based
database and simplifies the tasks of storing, retrieving, searching and sanitizing
favorite quotes..

## Goals
- Ease of storing favorite quotes
- Search with keywors and tags
- Profaniy checker
- Anonymous author

## 🚀 Quick Start
### Install 
- Clone the repo
- Run the project
```bash
go build -o quoify && ./quoify
```
- Open Localhost: 127.0.0.1:8080

## 📖 Usage

Available search methods:

* `Auhtor` - The author of quote
* `Tag` - The grouped quotes
* `Random` - Random quotes

## 🤝 Contributing

### Clone the repo

```bash
git clone https://github.com/xyz/zipzod@latest
cd zipzod
```

### Build the project

```bash
go build
```

### Run the project

```bash
./zipzod -i ./input -o ./output.zip
```

### Run the tests

```bash
go test ./...
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
