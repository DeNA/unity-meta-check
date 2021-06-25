package cli

type Command func(args []string, procInout ProcessInout, env Env) ExitStatus

type ExitStatus int

const (
	// ExitNormal means exit successfully.
	// SEE: http://tldp.org/LDP/abs/html/exitcodes.html
	ExitNormal   ExitStatus = 0

	// ExitAbnormal means exit not successfully.
	// SEE: http://tldp.org/LDP/abs/html/exitcodes.html
	ExitAbnormal ExitStatus = 1
)
