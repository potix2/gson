package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/potix2/gson"
)

func main() {
	bytes, err := ioutil.ReadFile("./example.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}

		ret, err := gson.Parse(bytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
		}
		fmt.Printf("%v", ret)
}

