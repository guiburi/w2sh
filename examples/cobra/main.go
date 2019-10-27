package main

import (
	"log"
	"net/http"

	"github.com/guiburi/w2sh"
	"github.com/guiburi/w2sh/examples/cobra/cmd"
)

func main() {
	http.HandleFunc("/", w2sh.Handle(cmd.RootCmd))
	log.Fatal(http.ListenAndServe(":8080", nil))
	//cmd.Execute()
}
