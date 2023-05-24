package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		env := make(Environment)
		env["FOO"] = EnvValue{Value: "foo"}
		env["BAR"] = EnvValue{NeedRemove: true}

		cmd := []string{"sh"}
		code := RunCmd(cmd, env)
		require.Equal(t, 0, code)
		require.Equal(t, "foo", os.Getenv("FOO"))
		envs := os.Environ()
		require.NotContains(t, envs, "BAR")
	})

	t.Run("unset variable", func(t *testing.T) {
		os.Setenv("FOO", "bar")
		env := Environment{
			"FOO": EnvValue{
				Value:      "foo",
				NeedRemove: true,
			},
		}
		code := RunCmd([]string{"sh"}, env)
		require.Equal(t, 0, code)
		envs := os.Environ()
		require.NotContains(t, envs, "FOO")
	})

	t.Run("nil env map", func(t *testing.T) {
		code := RunCmd([]string{"sh"}, nil)
		require.Equal(t, 0, code)
	})
}
