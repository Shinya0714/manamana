package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func main() {

	http.HandleFunc("/", echoHello)

	http.HandleFunc("/py", echoHello2)

	http.HandleFunc("/ls", echoHello3)
	// port
	http.ListenAndServe(":8000", nil)
}

func echoHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func echoHello2(w http.ResponseWriter, r *http.Request) {

	output, err := exec.Command("sh", "./hello.sh").Output()
	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Fprintf(w, string(output))
}

func echoHello3(w http.ResponseWriter, r *http.Request) {

	output, err := exec.Command("ls").Output()
	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Fprintf(w, string(output))
}
