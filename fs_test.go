package fs_test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-waitfor/waitfor"
	"github.com/go-waitfor/waitfor-fs"
)

func TestUse(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("waitfor_TestUse_%d.txt", time.Now().Nanosecond()))

	if err != nil {
		t.Error(err)
	}

	fileName := file.Name()

	defer file.Close()
	defer os.Remove(fileName)

	w := waitfor.New(fs.Use())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = w.Test(ctx, []string{"file://" + fileName})

	assert.NoError(t, err)
}

func TestFile(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("waitfor_TestFile_%d.txt", time.Now().Nanosecond()))

	assert.NoError(t, err)

	fileName := file.Name()

	defer file.Close()
	defer os.Remove(fileName)

	u, err := url.Parse("file://" + fileName)

	assert.NoError(t, err)

	r, err := fs.New(u)

	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = r.Test(ctx)

	assert.NoError(t, err)
}

func TestFile_FileNotExists(t *testing.T) {
	u, err := url.Parse("file://" + filepath.Join(os.TempDir(), "fdsfsdfds"))

	assert.NoError(t, err)

	r, err := fs.New(u)

	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = r.Test(ctx)

	assert.Error(t, err)
}

func TestFile_MissedURL(t *testing.T) {
	_, err := fs.New(nil)

	assert.Error(t, err)
}

func TestFile_ContextCanceled(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("waitfor_TestFile_Canceled_%d.txt", time.Now().Nanosecond()))

	assert.NoError(t, err)

	fileName := file.Name()

	defer file.Close()
	defer os.Remove(fileName)

	u, err := url.Parse("file://" + fileName)

	assert.NoError(t, err)

	r, err := fs.New(u)

	assert.NoError(t, err)

	// Create a canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err = r.Test(ctx)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestFile_ContextDeadlineExceeded(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("waitfor_TestFile_Deadline_%d.txt", time.Now().Nanosecond()))

	assert.NoError(t, err)

	fileName := file.Name()

	defer file.Close()
	defer os.Remove(fileName)

	u, err := url.Parse("file://" + fileName)

	assert.NoError(t, err)

	r, err := fs.New(u)

	assert.NoError(t, err)

	// Create a context that has already exceeded its deadline
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
	defer cancel()

	err = r.Test(ctx)

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestFile_EmptyPath(t *testing.T) {
	// Test with just "file://" which results in empty path
	u, err := url.Parse("file://")

	assert.NoError(t, err)

	r, err := fs.New(u)

	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = r.Test(ctx)

	assert.Error(t, err) // Empty path should not exist
}

func TestFile_NonFileScheme(t *testing.T) {
	// Test with http:// scheme - should still work but path will be malformed
	u, err := url.Parse("http://example.com/test")

	assert.NoError(t, err)

	r, err := fs.New(u)

	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = r.Test(ctx)

	assert.Error(t, err) // "http://example.com/test" is not a valid file path
}
