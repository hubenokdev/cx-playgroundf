package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/skycoin/cx-playground/playground"
	"github.com/skycoin/cx-playground/webapi"
	"github.com/skycoin/cx/cxparser/actions"
)

func main() {
	workingDir, _ := os.Getwd()
	if err := playground.InitPlayground(workingDir); err != nil {
		// error captured while initiating the playground examples, should be handled in the future
		fmt.Println("Fail to initiating palyground examples")
	}

	http.HandleFunc("/playground/examples", playground.GetExampleFileList)
	http.HandleFunc("/playground/examples/code", playground.GetExampleFileContent)

	http.Handle("/", http.FileServer(http.Dir("./dist")))
	http.Handle("/program/", webapi.NewAPI("/program", actions.AST))
	http.HandleFunc("/eval", playground.RunProgram)

	fmt.Println("Starting web service for CX playground on http://127.0.0.1:5336/")

	http.ListenAndServe(":5336", nil)
}
