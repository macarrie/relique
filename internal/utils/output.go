package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// OutputFormat holds the raw output flags of a CLI command. When none is set, the command displays
// its human readable output
type OutputFormat struct {
	JSON bool
	YAML bool
	TOML bool
}

func AddOutputFormatFlags(cobraCmd *cobra.Command, f *OutputFormat) {
	cobraCmd.Flags().BoolVarP(&f.JSON, "json", "", false, "Output raw data as JSON")
	cobraCmd.Flags().BoolVarP(&f.YAML, "yaml", "", false, "Output raw data as YAML")
	cobraCmd.Flags().BoolVarP(&f.TOML, "toml", "", false, "Output raw data as TOML")
	cobraCmd.MarkFlagsMutuallyExclusive("json", "yaml", "toml")
}

// IsRaw returns true when a raw output format has been requested
func (f OutputFormat) IsRaw() bool {
	return f.JSON || f.YAML || f.TOML
}

// WriteRaw serializes data in the requested raw format.
// JSON and YAML use the same snake_case field names as the HTTP API
func (f OutputFormat) WriteRaw(w io.Writer, data interface{}) error {
	var out []byte
	var err error

	switch {
	case f.JSON:
		out, err = json.MarshalIndent(data, "", "  ")
		out = append(out, '\n')
	case f.YAML:
		out, err = toYAML(data)
	case f.TOML:
		out, err = toml.Marshal(data)
	default:
		return fmt.Errorf("no raw output format requested")
	}
	if err != nil {
		return fmt.Errorf("cannot serialize data: %w", err)
	}

	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("cannot write output: %w", err)
	}

	return nil
}

// toYAML goes through JSON so that YAML keys follow the json struct tags instead of
// yaml.v3 lowercased field names
func toYAML(data interface{}) ([]byte, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize data to json: %w", err)
	}

	var generic interface{}
	if err := json.Unmarshal(raw, &generic); err != nil {
		return nil, fmt.Errorf("cannot convert json data: %w", err)
	}

	return yaml.Marshal(generic)
}
