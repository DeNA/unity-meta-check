package yaml

import (
	"testing"
)

var _ InputDef = StringInputDef{}
var _ InputDef = BoolInputDef{}
var _ InputDef = IntInputDef{}

func TestInputDefs(t *testing.T) {
	for i, inputDef := range BuildMetadata().Inputs {
		if inputDef.Name() == "" {
			t.Errorf("InputDefs[%d] must have a non-empty name", i)
		}

		if inputDef.Desc() == "" {
			t.Errorf("%q must have a non-empty description", inputDef.Name())
		}

		if !inputDef.Required() && inputDef.DefaultValueAsYAML() == "" {
			t.Errorf("optional input %q must have a non-empty default value representation for YAML", inputDef.Name())
		}
	}
}
