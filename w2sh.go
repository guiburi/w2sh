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
	"sync"
)

type Page struct {
	Title      string
	Output     string
	Form       string
	CommandMap map[string]*cobra.Command
}

var mux sync.Mutex

func Handle(root *cobra.Command) func(w http.ResponseWriter, r *http.Request) {
	//create cmd map
	cmds := make(map[string]*cobra.Command)
	cmds[root.Name()] = root
	for _, c := range root.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		cmds[c.Name()] = c
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//todo: lookup command
		var cmd *cobra.Command
		key := r.URL.Path[1:len(r.URL.Path)]
		cmd, ok := cmds[key]
		if !ok || r.URL.Path == "/" {
			cmd = root
		}

		var output *Page
		t := strings.Title(strings.ToLower(cmd.Name()))
		f := genForm(cmd)

		args := []string{}
		if r.URL.Path[1:len(r.URL.Path)] != root.Name() {
			args = append(args, r.URL.Path[1:len(r.URL.Path)])
		}

		switch r.Method {
		case http.MethodGet:
			// todo: help doesnt exit. why??
			//args = append(args, "--help")
			//o, _ := executeCommand(root, args...)
			//output = &Page{t, o, f, cmds}
			output = &Page{t, cmd.UsageString(), f, cmds}

		case http.MethodPost:
			_ = r.ParseForm()
			for k, v := range r.Form {
				args = append(args, "--"+k)
				args = append(args, strings.Join(v, ""))
			}

			// todo: remove this
			fmt.Print(args)
			fmt.Print("\n")

			//execute
			o, _ := executeCommand(root, args...)
			output = &Page{t, o, f, cmds}
		default:
			// todo: give an error message.
		}
		// render
		_ = tmpl.Execute(w, output)
	}
}

func genForm(cmd *cobra.Command) (form string) {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf(`<form method="POST" action="/%s">`, cmd.Name()))

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

func executeCommand(cmd *cobra.Command, args ...string) (output string, err error) {
	mux.Lock()
	defer mux.Unlock()

	done := capture()
	cmd.SetArgs(args)
	_, err = cmd.ExecuteC()
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
	saveErr := os.Stderr
	os.Stdout = w
	os.Stderr = w
	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		err = r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		os.Stderr = saveErr
		err := w.Close()
		err = <-done
		return buf.String(), err
	}
}
