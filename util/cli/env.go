package cli

import "os"

type Env func(key string) string

func NewEnv() Env {
	return func(key string) string {
		return os.Getenv(key)
	}
}
