package cache

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/restic/restic/internal/backend"
	"github.com/restic/restic/internal/restic"
	rtest "github.com/restic/restic/internal/test"
)

func TestNew(t *testing.T) {
	parent := rtest.TempDir(t)
	basedir := filepath.Join(parent, "cache")
	id := restic.NewRandomID().String()
	tagFile := filepath.Join(basedir, "CACHEDIR.TAG")
	versionFile := filepath.Join(basedir, id, "version")

	const (
		stepCreate = iota
		stepComplete
		stepRmTag
		stepRmVersion
		stepEnd
	)

	for step := stepCreate; step < stepEnd; step++ {
		switch step {
		case stepRmTag:
			rtest.OK(t, os.Remove(tagFile))
		case stepRmVersion:
			rtest.OK(t, os.Remove(versionFile))
		}

		c, err := New(id, basedir, 0)
		rtest.OK(t, err)
		rtest.Equals(t, basedir, c.Base)
		rtest.Equals(t, step == stepCreate, c.Created)

		for _, name := range []string{tagFile, versionFile} {
			info, err := os.Lstat(name)
			rtest.OK(t, err)
			rtest.Assert(t, info.Mode().IsRegular(), "")
		}
	}
}

// TestCacheSize tests that the cache respects the maximum size setting.
func TestCacheSize(t *testing.T) {
	tempdir := rtest.TempDir(t)

	// Create a cache with a 100 byte size limit
	const maxSize int64 = 100
	c, err := New("test-repo", tempdir, maxSize)
	if err != nil {
		t.Fatal(err)
	}

	// Create test content larger than the maximum cache size
	content := []byte("test") // 4 bytes per file

	// Save several files to the cache to exceed the limit
	for i := range make([]int, 10) {
		h := backend.Handle{Type: restic.PackFile, Name: fmt.Sprintf("test-file-%d", i)}
		err = c.save(h, bytes.NewReader(content))
		if err != nil {
			t.Fatal(err)
		}
	}

	// Get the current cache size and verify it's under the limit
	size, err := c.getTotalSize()
	if err != nil {
		t.Fatal(err)
	}

	// Allow for small variation due to cleanup targeting 80% of max
	if size > maxSize {
		t.Errorf("cache size %d exceeds maximum size %d", size, maxSize)
	}

	// Check that the oldest files were removed
	for i := 0; i < 5; i++ {
		h := backend.Handle{Type: restic.PackFile, Name: fmt.Sprintf("test-file-%d", i)}
		has := c.Has(h)
		// The first few files should have been removed
		if has {
			t.Errorf("expected file %v to be removed from cache, but it still exists", h)
		}
	}
}
