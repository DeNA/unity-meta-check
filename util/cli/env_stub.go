package cli

func AnyEnv() Env {
	return StubEnv("ANY_ENV")
}

func StubEnv(result string) Env {
	return func(string) string {
		return result
	}
}
