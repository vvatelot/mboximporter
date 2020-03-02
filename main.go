package main

import (
	"fmt"

	"github.com/jaffee/commandeer"

	"github.com/rpsl/mboximporter/cmd"
)

func main() {
	err := commandeer.Run(cmd.NewImport())

	if err != nil {
		fmt.Println(err)
	}
}
