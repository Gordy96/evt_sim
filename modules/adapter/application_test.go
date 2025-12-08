package adapter

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
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

	l := zaptest.NewLogger(t)

	a, err := New(
		lib,
		WithParam[int]("counter", 3),
		WithParam[string]("name", "foobar"),
		WithParam[float64]("factor", 12.34),
		WithLogger(func(i int, s string) {
			l.Log(zapcore.Level(i), s)
		}),
	)
	assert.NoError(t, err)

	assert.NoError(t, a.Init(nil, &port))
	defer a.Close()

	for i := 0; i < 3; i++ {
		err = a.TriggerPortInterrupt("port")
		assert.NoError(t, err)
	}

	assert.ElementsMatch(
		t,
		[]string{
			"foobar 12.340000 hello world 2",
			"foobar 12.340000 hello world 1",
			"foobar 12.340000 hello world 0",
		},
		port.recorded,
	)
}
