package memory

import (
	"code.jogchat.internal/go-schemaless/storagetest"
	"testing"
)

func TestMemory(t *testing.T) {
	m := New()
	storagetest.StorageTest(t, m)
}
