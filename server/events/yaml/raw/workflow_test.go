package raw_test

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/runatlantis/atlantis/server/events/yaml/raw"
	"github.com/runatlantis/atlantis/server/events/yaml/valid"
	. "github.com/runatlantis/atlantis/testing"
	"gopkg.in/yaml.v2"
)

func TestWorkflow_UnmarshalYAML(t *testing.T) {
	cases := []struct {
		description string
		input       string
		exp         raw.Workflow
		expErr      string
	}{
		{
			description: "empty",
			input:       ``,
			exp: raw.Workflow{
				Apply: nil,
				Plan:  nil,
			},
		},
		{
			description: "yaml null",
			input:       `~`,
			exp: raw.Workflow{
				Apply: nil,
				Plan:  nil,
			},
		},
		{
			description: "only plan/apply set",
			input: `
plan:
apply:
`,
			exp: raw.Workflow{
				Apply: nil,
				Plan:  nil,
			},
		},
		{
			description: "steps set to null",
			input: `
plan:
  steps: ~
apply:
  steps: ~`,
			exp: raw.Workflow{
				Plan: &raw.Stage{
					Steps: nil,
				},
				Apply: &raw.Stage{
					Steps: nil,
				},
			},
		},
		{
			description: "steps set to empty slice",
			input: `
plan:
  steps: []
apply:
  steps: []`,
			exp: raw.Workflow{
				Plan: &raw.Stage{
					Steps: []raw.Step{},
				},
				Apply: &raw.Stage{
					Steps: []raw.Step{},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			var w raw.Workflow
			err := yaml.UnmarshalStrict([]byte(c.input), &w)
			if c.expErr != "" {
				ErrEquals(t, c.expErr, err)
				return
			}
			Ok(t, err)
			Equals(t, c.exp, w)
		})
	}
}

func TestWorkflow_Validate(t *testing.T) {
	// Should call the validate of Stage.
	w := raw.Workflow{
		Apply: &raw.Stage{
			Steps: []raw.Step{
				{
					Key: String("invalid"),
				},
			},
		},
	}
	validation.ErrorTag = "yaml"
	ErrEquals(t, "apply: (steps: (0: \"invalid\" is not a valid step type.).).", w.Validate())

	// Unset keys should validate.
	Ok(t, (raw.Workflow{}).Validate())
}

func TestWorkflow_ToValid(t *testing.T) {
	cases := []struct {
		description string
		input       raw.Workflow
		exp         valid.Workflow
	}{
		{
			description: "nothing set",
			input:       raw.Workflow{},
			exp: valid.Workflow{
				Apply: nil,
				Plan:  nil,
			},
		},
		{
			description: "fields set",
			input: raw.Workflow{
				Apply: &raw.Stage{
					Steps: []raw.Step{
						{
							Key: String("init"),
						},
					},
				},
				Plan: &raw.Stage{
					Steps: []raw.Step{
						{
							Key: String("init"),
						},
					},
				},
			},
			exp: valid.Workflow{
				Apply: &valid.Stage{
					Steps: []valid.Step{
						{
							StepName: "init",
						},
					},
				},
				Plan: &valid.Stage{
					Steps: []valid.Step{
						{
							StepName: "init",
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			Equals(t, c.exp, c.input.ToValid())
		})
	}
}
