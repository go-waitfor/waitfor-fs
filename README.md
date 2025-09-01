# waitfor-fs

[![Go Version](https://img.shields.io/github/go-mod/go-version/go-waitfor/waitfor-fs)](https://golang.org/dl/)
[![Build Status](https://github.com/go-waitfor/waitfor-fs/workflows/Build/badge.svg)](https://github.com/go-waitfor/waitfor-fs/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-waitfor/waitfor-fs)](https://goreportcard.com/report/github.com/go-waitfor/waitfor-fs)

File system resource readiness assertion library for Go. This library extends the [waitfor](https://github.com/go-waitfor/waitfor) framework to provide file system checking capabilities, allowing you to wait for files and directories to exist before proceeding with your application logic.

## Features

- **File existence checking**: Wait for files to be created or become available
- **Directory existence checking**: Wait for directories to be created
- **Context support**: Full support for Go contexts with cancellation and timeouts
- **Configurable retries**: Use the waitfor framework's built-in retry mechanisms
- **Thread-safe**: Safe for concurrent use

## Installation

```bash
go get github.com/go-waitfor/waitfor-fs
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-waitfor/waitfor"
	fs "github.com/go-waitfor/waitfor-fs"
)

func main() {
	// Create a waitfor runner with file system support
	runner := waitfor.New(fs.Use())

	// Wait for files and directories to exist
	err := runner.Test(
		context.Background(),
		[]string{"file://./my-file.txt", "file://./my-dir/"},
		waitfor.WithAttempts(5),
	)

	if err != nil {
		fmt.Printf("Files not ready: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All files are ready!")
}
```

## URL Format

The library uses the `file://` URL scheme to specify file system resources:

- **Files**: `file:///absolute/path/to/file.txt` or `file://./relative/path/file.txt`
- **Directories**: `file:///absolute/path/to/directory/` or `file://./relative/directory/`

## Examples

### Basic File Checking

```go
package main

import (
	"context"
	"time"

	"github.com/go-waitfor/waitfor"
	fs "github.com/go-waitfor/waitfor-fs"
)

func main() {
	runner := waitfor.New(fs.Use())
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Wait for a configuration file to be created
	err := runner.Test(ctx, []string{"file://./config.json"})
	if err != nil {
		panic(err)
	}
}
```

### Multiple Files with Custom Configuration

```go
package main

import (
	"context"
	"time"

	"github.com/go-waitfor/waitfor"
	fs "github.com/go-waitfor/waitfor-fs"
)

func main() {
	runner := waitfor.New(fs.Use())
	
	ctx := context.Background()
	
	// Wait for multiple files with custom retry configuration
	err := runner.Test(
		ctx,
		[]string{
			"file://./data/input.csv",
			"file://./data/schema.json",
			"file://./logs/",
		},
		waitfor.WithAttempts(10),
		waitfor.WithDelay(2*time.Second),
	)
	
	if err != nil {
		panic(err)
	}
}
```

### Using the Resource Directly

```go
package main

import (
	"context"
	"net/url"

	fs "github.com/go-waitfor/waitfor-fs"
)

func main() {
	// Parse the file URL
	u, err := url.Parse("file://./important-file.txt")
	if err != nil {
		panic(err)
	}

	// Create the file resource
	resource, err := fs.New(u)
	if err != nil {
		panic(err)
	}

	// Test if the file exists
	ctx := context.Background()
	err = resource.Test(ctx)
	if err != nil {
		panic(err)
	}
}
```

## API Reference

### `func Use() waitfor.ResourceConfig`

Returns a resource configuration for the waitfor framework that enables file system checking with the `file` scheme.

### `func New(u *url.URL) (waitfor.Resource, error)`

Creates a new file system resource from a URL. The URL must use the `file` scheme.

**Parameters:**
- `u`: A URL with the `file` scheme pointing to a file or directory

**Returns:**
- A `waitfor.Resource` that can be used to test file existence
- An error if the URL is invalid or nil

### `func (*File) Test(ctx context.Context) error`

Tests whether the file or directory exists. This method is called by the waitfor framework during resource checking.

**Parameters:**
- `ctx`: Context for cancellation and timeout control

**Returns:**
- `nil` if the file or directory exists
- An error if the file doesn't exist, the context is cancelled, or there's a file system error

## About waitfor

This library is part of the [waitfor](https://github.com/go-waitfor/waitfor) ecosystem, which provides a unified framework for waiting on various types of resources to become ready. Other available waitfor plugins include:

- HTTP endpoints
- Database connections  
- Network services
- And more...

## Requirements

- Go 1.23 or later

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
