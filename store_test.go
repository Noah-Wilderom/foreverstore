package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "testing-transform-func"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "e6214631ee5b5b936bd57dfeb3ba174a756b9cbc"
	expectedPathName := "e6214/631ee/5b5b9/36bd5/7dfeb/3ba17/4a756/b9cbc"

	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s want %s\n", pathKey.PathName, expectedPathName)
	}

	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("have %s want %s\n", pathKey.Filename, expectedOriginalKey)
	}
}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	key := "testing-key"

	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	key := "testing-key"

	defer s.Delete(key)
	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	if string(b) != string(data) {
		t.Errorf("want %s have %s", string(data), string(b))
	}
}
