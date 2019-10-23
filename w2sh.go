package w2sh

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Title  string
	Output string
	Form   string
}

func Handle(root *cobra.Command) func(w http.ResponseWriter, r *http.Request) {
	// todo: create map of commands and sub commands
	return func(w http.ResponseWriter, r *http.Request) {
		//todo: lookup command
		var output *Page
		switch r.Method {
		case http.MethodGet:
			t := strings.Title(strings.ToLower(root.Name()))
			f := genForm(root)
			output = &Page{t, root.UsageString(), f}
		case http.MethodPost:
			//todo: extract flags
			_ = r.ParseForm()
			name := r.Form.Get("name")
			t := strings.Title(strings.ToLower(root.Name()))
			f := genForm(root)
			//execute
			o, _ := executeCommand(root, "--name", name)
			output = &Page{t, o, f}
		default:
			// todo: give an error message.
		}
		// render
		_ = tmpl.Execute(w, output)
	}
}

func genForm(cmd *cobra.Command) (form string) {
	buf := new(bytes.Buffer)
	buf.WriteString(`<form method="POST">`)

	flags := cmd.NonInheritedFlags()
	if flags.HasAvailableFlags() {
		genInput(buf, flags)
	}

	flags = cmd.InheritedFlags()
	if flags.HasAvailableFlags() {
		genInput(buf, flags)
	}
	buf.WriteString(`</br><input type="submit" value="Submit">`)
	buf.WriteString(`</form>`)
	return buf.String()
}

func genInput(buf *bytes.Buffer, flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		if len(flag.Deprecated) > 0 || flag.Hidden || flag.Name == "help" {
			return
		}
		format := fmt.Sprintf(`%s: <input type="text" name="%s"><br>`, flag.Name, flag.Name)
		buf.WriteString(fmt.Sprintf(format))
	})
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	done := capture()
	root.SetArgs(args)
	_, err = root.ExecuteC()
	if err != nil {
		return "", err
	}
	output, err = done()
	return output, err
}

func capture() func() (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return func() (string, error) {
			return "", err
		}
	}
	done := make(chan error, 1)
	save := os.Stdout
	os.Stdout = w
	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		err = r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		err := w.Close()
		err = <-done
		return buf.String(), err
	}
}
