package faeio

import (
	"io"
	"os"
)

// MysticalPortal represents the magical connection between the forest and the outside world
var MysticalPortal = &struct {
	io.ReadCloser
	io.Writer
}{
	io.NopCloser(os.Stdin),
	os.Stdout,
}
