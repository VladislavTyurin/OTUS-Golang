package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("invalid name", func(t *testing.T) {
		f, err := os.Create(path.Join("testdata", "FOO="))
		require.NoError(t, err)
		defer os.Remove(f.Name())
		env, err := ReadDir("testdata")
		require.ErrorIs(t, err, ErrInvalidName)
		require.Nil(t, env)
	})

	t.Run("not a dir error", func(t *testing.T) {
		_, err := ReadDir(path.Join("testdata", "echo.sh"))
		require.Error(t, err)
	})

	t.Run("empty first line", func(t *testing.T) {
		f, err := os.Create(path.Join("testdata", "MYFOO"))
		require.NoError(t, err)
		defer os.Remove(f.Name())

		f.WriteString("\nvalue\n")
		f.Close()
		env, err := ReadDir("testdata")
		require.NoError(t, err)
		require.Equal(t, env["FOO"].Value, "")
	})
}
