package gophash

import (
	"encoding/hex"
	"log"
	"os"
	"path"
	"testing"
)

const testdata = "testdata"

// ffmpeg -i test.webp -compression_level 0 -frames:v 1 test.png
func TestPHash(t *testing.T) {
	phash := New(path.Join(testdata, "test.png"), log.New(os.Stderr, "[phash]", log.LstdFlags))
	t.Log(hex.EncodeToString(phash.Sum(nil)))
}
