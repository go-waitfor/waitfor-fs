package main

import (
"fmt"
"net/url"
"context"
"time"
"os"
fs "github.com/go-waitfor/waitfor-fs"
)

func main() {
// Test various URL scenarios
testCases := []string{
"file://",
"file:///",
"http://example.com",
"file://./test",
"file://test with spaces",
}

for _, tc := range testCases {
fmt.Printf("Testing URL: %q\n", tc)
u, err := url.Parse(tc)
if err != nil {
fmt.Printf("  Parse error: %v\n", err)
continue
}

r, err := fs.New(u)
if err != nil {
fmt.Printf("  New error: %v\n", err)
continue
}

fmt.Printf("  URL String: %q\n", u.String())
fmt.Printf("  Path after TrimPrefix: %q\n", u.String()[len("file://"):])

ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
err = r.Test(ctx)
cancel()

if err != nil {
fmt.Printf("  Test error: %v\n", err)
} else {
fmt.Printf("  Test: OK\n")
}
fmt.Println()
}
}
