package yaml

import (
	"fmt"
)

type (
	Name string
	Desc string

	InputDef interface {
		Name() Name
		Desc() Desc
		Required() bool
		DefaultValueAsYAML() string
	}
)

type StringInputDef struct {
	name         Name
	desc         Desc
	required     bool
	defaultValue string
}

func (s StringInputDef) Name() Name {
	return s.name
}

func (s StringInputDef) Desc() Desc {
	return s.desc
}

func (s StringInputDef) Required() bool {
	return s.required
}

func (s StringInputDef) DefaultValueAsYAML() string {
	return fmt.Sprintf("%q", s.defaultValue)
}

type BoolInputDef struct {
	name         Name
	desc         Desc
	required     bool
	defaultValue bool
}

func (b BoolInputDef) Name() Name {
	return b.name
}

func (b BoolInputDef) Desc() Desc {
	return b.desc
}

func (b BoolInputDef) Required() bool {
	return b.required
}

func (b BoolInputDef) DefaultValueAsYAML() string {
	return fmt.Sprintf("%t", b.defaultValue)
}

type IntInputDef struct {
	name         Name
	desc         Desc
	required     bool
	defaultValue int
}

func (u IntInputDef) Name() Name {
	return u.name
}

func (u IntInputDef) Desc() Desc {
	return u.desc
}

func (u IntInputDef) Required() bool {
	return u.required
}

func (u IntInputDef) DefaultValueAsYAML() string {
	return fmt.Sprintf("%d", u.defaultValue)
}
