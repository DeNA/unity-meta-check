package cli

func AnyEnv() Env {
	return ConstEnv("ANY_ENV")
}

func ConstEnv(result string) Env {
	return func(string) string {
		return result
	}
}

func StubEnv(m map[string]string) Env {
	return func(n string) string {
		return m[n]
	}
}