package main

import (
	"github.com/guiburi/w2sh"
	"github.com/guiburi/w2sh/examples/cobra/cmd"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/",w2sh.Handle(cmd.GetRoot()))
	log.Fatal(http.ListenAndServe(":8080", nil))
	//cmd.Execute()
}
