package adapter

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func compileSo(ctx context.Context, path string) error {
	return exec.CommandContext(
		ctx,
		"gcc",
		"-fPIC",
		"-shared",
		"-o",
		filepath.Join(path, "testdata", "plugin.so"),
		filepath.Join(path, "testdata", "plugin.c"),
	).Run()
}

type FakePort struct {
	recorded []string
}

func (f *FakePort) Name() string {
	return "port"
}

func (f *FakePort) Read(b []byte) (n int, err error) {
	const hw = "hello world"
	copy(b, hw)

	return len(hw), nil
}

func (f *FakePort) Write(b []byte) (n int, err error) {
	f.recorded = append(f.recorded, string(b))
	return len(b), nil
}

func TestCompile(t *testing.T) {
	path, _ := os.Getwd()
	err := compileSo(t.Context(), path)
	assert.NoError(t, err)

	lib, err := OpenLib(filepath.Join(path, "testdata", "plugin.so"))
	assert.NoError(t, err)
	defer lib.Release()

	port := FakePort{}

	a, err := New("runner", []Port{&port}, nil, lib)
	assert.NoError(t, err)

	a.Init()
	defer a.Close()

	for i := 0; i < 3; i++ {
		err = a.TriggerPinInterrupt(2)
		assert.NoError(t, err)
	}

	assert.ElementsMatch(
		t,
		[]string{
			"hello world 1",
			"hello world 2",
			"hello world 3",
		},
		port.recorded,
	)
}
