// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(http.Dir("raw"), vfsgen.Options{
		PackageName:  "static",
		VariableName: "Content",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
