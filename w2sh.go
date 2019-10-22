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
	"text/template"
)

type Page struct {
	Output string
	Form   string
}

var counter int

var t = template.Must(template.ParseFiles("/Users/ali.firash/code/w2sh/template.html"))

func Handle(root *cobra.Command) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		counter++
		//find correct command
		//extract flags from cmd
		//execute with flags if method get run help
		var output Page
		if r.Method == "POST" {
			fmt.Printf("%d::POSTED\n", counter)
			output, _ = html(root, "--name", "seek")
		} else {
			fmt.Printf("%d::GETTED\n", counter)
			output, _ = html(root, "--help")
		}
		t.Execute(w, output)
	}
}

func html(cmd *cobra.Command, args ...string) (page Page, err error) {
	form := genForm(cmd)
	output, err := executeCommand(cmd, args...)
	if err != nil {
		return Page{Output: err.Error()}, err
	}
	return Page{Output: output, Form: form}, nil
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
	buf.WriteString(`<input type="submit" value="Submit">`)
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
		r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		w.Close()
		err := <-done
		return buf.String(), err
	}
}
