# waitfor-fs
File system resource readiness assertion library

# Quick start

```go
package main

import (
	"context"
	"fmt"
	"github.com/go-waitfor/waitfor"
	"github.com/go-waitfor/waitfor-fs"
	"os"
)

func main() {
	runner := waitfor.New(fs.Use())

	err := runner.Test(
		context.Background(),
		[]string{"fs://./my-file.txt", "fs://./my-dir/"},
		waitfor.WithAttempts(5),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```
