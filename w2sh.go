package w2sh

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandMap struct {
	Cmds     map[string]*cobra.Command
	Children []*CommandMap
}

type Page struct {
	Title  string
	Output string
	Form   string
	CmdMap CommandMap
}

var mux sync.Mutex

func Handle(root *cobra.Command) func(w http.ResponseWriter, r *http.Request) {
	//create cmd map
	var cmdMap = &CommandMap{}
	cmdMap.Cmds = make(map[string]*cobra.Command)
	cmdMap = extract(root, cmdMap)

	return func(w http.ResponseWriter, r *http.Request) {
		args := []string{}

		args, cmd := lookup(r.URL.Path[1:len(r.URL.Path)], cmdMap, args)
		if cmd == nil || r.URL.Path == "/" {
			fmt.Print("was not found \n")
			cmd = root
		}

		var output *Page
		t := strings.Title(strings.ToLower(cmd.Name()))
		f := genForm(cmd)

		//if r.URL.Path[1:len(r.URL.Path)] != root.Name() {
		//	args = append(args, r.URL.Path[1:len(r.URL.Path)])
		//}

		switch r.Method {
		case http.MethodGet:
			// todo: help doesnt exit. why??
			//args = append(args, "--help")
			//o, _ := executeCommand(root, args...)
			//output = &Page{t, o, f, cmds}
			output = &Page{t, cmd.UsageString(), f, *cmdMap}

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
			o, _ := execute(root, args...)
			output = &Page{t, o, f, *cmdMap}
		default:
			// todo: give an error message.
		}
		// render
		_ = tmpl.Execute(w, output)
	}
}

func extract(cmd *cobra.Command, cmdMap *CommandMap) *CommandMap {
	fmt.Print("extracting::" + cmd.Name() + "\n")
	cmdMap.Cmds[cmd.Name()] = cmd

	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		var childMap = &CommandMap{}
		childMap.Cmds = make(map[string]*cobra.Command)
		childMap.Cmds[c.Name()] = c
		cmdMap.Children = append(cmdMap.Children, childMap)
		extract(c, childMap)
	}
	return cmdMap
}

func lookup(name string, cmdMap *CommandMap, args []string) ([]string, *cobra.Command) {

	if cmd, ok := cmdMap.Cmds[name]; ok {
		if cmd.HasParent() {
			args = append(args, cmd.Name())
		}
		fmt.Print("root::" + cmd.Name() + "\n")
		return args, cmd
	}

	for _, child := range cmdMap.Children {
		if cmd, ok := child.Cmds[name]; ok {
			args = append(args, cmd.Name())
			fmt.Printf("args::%v\n", args)
			return args, cmd
		}

		lookup(name, child, args)
	}
	return nil, nil
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

func execute(cmd *cobra.Command, args ...string) (output string, err error) {
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
