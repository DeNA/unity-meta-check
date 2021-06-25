package opt

import (
	"fmt"
	"sort"
)

type (
	Type string
	Name string
	Desc string

	Option interface {
		Required() bool
		Name() Name
		Desc() Desc
		Type() Type
		DefaultValueText() string
	}
)

const (
	TypeString Type = "string"
	TypeBool   Type = "bool"
)

func NewOptionalStringOption(name Name, desc Desc, defaultValue string) StringOption {
	return StringOption{
		required:     false,
		name:         name,
		desc:         desc,
		DefaultValue: defaultValue,
	}
}

func NewRequiredStringOption(name Name, desc Desc) StringOption {
	return StringOption{
		required: true,
		name:     name,
		desc:     desc,
	}
}

type StringOption struct {
	required     bool
	name         Name
	desc         Desc
	DefaultValue string
}

func (s StringOption) Name() Name {
	return s.name
}

func (s StringOption) Desc() Desc {
	return s.desc
}

func (s StringOption) Type() Type {
	return TypeString
}

func (s StringOption) Required() bool {
	return s.required
}

func (s StringOption) DefaultValueText() string {
	return fmt.Sprintf("%q", s.DefaultValue)
}

func NewOptionalBoolOption(name Name, desc Desc, defaultValue bool) BoolOption {
	return BoolOption{
		required:     false,
		name:         name,
		desc:         desc,
		DefaultValue: defaultValue,
	}
}

func NewRequiredBoolOption(name Name, desc Desc) BoolOption {
	return BoolOption{
		required: true,
		name:     name,
		desc:     desc,
	}
}

type BoolOption struct {
	required     bool
	name         Name
	desc         Desc
	DefaultValue bool
}

func (b BoolOption) Name() Name {
	return b.name
}

func (b BoolOption) Desc() Desc {
	return b.desc
}

func (b BoolOption) Type() Type {
	return TypeBool
}

func (b BoolOption) Required() bool {
	return b.required
}

func (b BoolOption) DefaultValueText() string {
	return fmt.Sprintf("%t", b.DefaultValue)
}

func Sort(flags []Option) {
	sort.Slice(flags, func(i, j int) bool {
		if flags[i].Required() != flags[i].Required() {
			return flags[i].Required()
		}
		return flags[i].Name() > flags[j].Name()
	})
}
