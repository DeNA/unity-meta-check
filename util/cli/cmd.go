package cli

type Command func(args []string, procInout ProcessInout, env Env) ExitStatus

type ExitStatus int

const (
	// http://tldp.org/LDP/abs/html/exitcodes.html
	ExitNormal   ExitStatus = 0
	ExitAbnormal ExitStatus = 1
)
