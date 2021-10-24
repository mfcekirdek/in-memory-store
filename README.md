# in-memory-store [![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=Check%20out%20this%20cool%20project&url=https://gitlab.com/mfcekirdek/in-memory-store&hashtags=project,opensource)
![Github License](https://img.shields.io/badge/license-MIT-green)
![Go Version](https://img.shields.io/badge/go-1.16-red.svg)

## Description
A simple in memory key value store application.

## Table of content
- [Project setup](#project-setup)
- [Run lint checks](#run-lint-checks)
- [Run Tests](#run-tests)
- [Build project](#build-project)
- [Run](#run)
- [Built With](#built-with)

### Project setup
``` console
go mod download
```

### Run lint checks
``` console
make lint
```

### RUN tests
``` console
make test
```

### Build project
``` console
make build
```

### Run
You can simply run the app with default configurations:
``` console
make run
```

Or set configs by exporting environment variables.
``` console
export IS_DEBUG=true
export APP_NAME=kv-store
export PORT=8080
export STORAGE_DIR_PATH=storage
export SAVE_TO_FILE_INTERVAL=10
make run
```

### Built With
1. Go
    - Go standard library
    - github.com/spf13/viper --> For reading env configs at startup.
    - github.com/k0kubun/pp --> For printing app configs in a pretty format.
    - github.com/golang/mock --> For mocking components in tests.