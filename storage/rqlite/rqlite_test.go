package rqlite

import (
	"code.jogchat.internal/go-schemaless/storagetest"
	"testing"
)

func TestRQLite(t *testing.T) {
	m := New().WithZap().WithURL("http://")
	storagetest.StorageTest(t, m)
}
