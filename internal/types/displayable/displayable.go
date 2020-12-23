package displayable

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pkg/errors"
)

const (
	TUI = iota
	JSON
)

var DisplayMode int

type Displayable interface {
	Display() Struct
}

type Struct interface {
	Summary() string
	Details() string
	TableRow() []string
	TableHeaders() []string
}

func Table(table []Displayable) {
	var out []Struct
	for _, d := range table {
		out = append(out, d.Display())
	}

	if DisplayMode == JSON {
		j, err := toJson(out)
		if err != nil {
			fmt.Printf("Cannot display item in json format: '%s'\n", err)
			return
		}

		fmt.Printf("%v\n", j)
		return
	} else {
		if len(out) == 0 {
			fmt.Println("No elements to display")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 4, 5, ' ', 0)
		defer w.Flush()

		headers := out[0].TableHeaders()

		for _, s := range headers {
			fmt.Fprintf(w, "%s\t", s)
		}
		fmt.Fprintln(w, "")

		for _, item := range out {
			row := item.TableRow()
			for _, s := range row {
				fmt.Fprintf(w, "%s\t", s)
			}
			fmt.Fprintln(w, "")
		}
		fmt.Fprintf(w, "\n%d elements displayed\n", len(out))
	}
}

func Details(d Displayable) {
	out := d.Display()

	if DisplayMode == JSON {
		j, err := toJson(out)
		if err != nil {
			fmt.Printf("Cannot display item in json format: '%s'\n", err)
			return
		}

		fmt.Printf("%v\n", j)
		return
	} else {
		// TODO: Pretty text display
		fmt.Printf("%v\n", out.Details())
	}
}

func toJson(item interface{}) (string, error) {
	out, err := json.Marshal(item)
	if err != nil {
		return "", errors.Wrap(err, "cannot convert displayable item to json")
	}

	return string(out), nil
}
