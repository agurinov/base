package tools

import (
	"bytes"
	"reflect"
	"testing"
)

func TestToCloser(t *testing.T) {
	noCloser := bytes.NewBuffer([]byte{})
	nativeReadCloser := ReadCloser(noCloser)
	nativeWriteCloser := WriteCloser(noCloser)

	t.Run("Read", func(t *testing.T) {
		t.Run("native", func(t *testing.T) {
			oldPtr := reflect.ValueOf(nativeReadCloser).Pointer()
			newPtr := reflect.ValueOf(ReadCloser(nativeReadCloser)).Pointer()

			// native ReadCloser will return without any injections -> same pointer
			if oldPtr != newPtr {
				t.Error("Unexpected pointer")
			}
		})

		t.Run("obtained", func(t *testing.T) {
			oldPtr := reflect.ValueOf(noCloser).Pointer()
			newPtr := reflect.ValueOf(ReadCloser(noCloser)).Pointer()

			// no ReadCloser will be returned with injection of .close() method (just return nil)
			// -> different pointer
			if oldPtr == newPtr {
				t.Error("Unexpected pointer")
			}
		})
	})

	t.Run("Write", func(t *testing.T) {
		t.Run("native", func(t *testing.T) {
			oldPtr := reflect.ValueOf(nativeWriteCloser).Pointer()
			newPtr := reflect.ValueOf(WriteCloser(nativeWriteCloser)).Pointer()

			// native WriteCloser will return without any injections -> same pointer
			if oldPtr != newPtr {
				t.Error("Unexpected pointer")
			}
		})

		t.Run("obtained", func(t *testing.T) {
			oldPtr := reflect.ValueOf(noCloser).Pointer()
			newPtr := reflect.ValueOf(WriteCloser(noCloser)).Pointer()

			// no WriteCloser will be returned with injection of .close() method (just return nil)
			// -> different pointer
			if oldPtr == newPtr {
				t.Error("Unexpected pointer")
			}
		})
	})
}
