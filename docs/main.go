package main

import "github.com/spf13/cobra/doc"
import "github.com/apokalyptik/pgm/cmd"

func main() {
	doc.GenMarkdownTree(cmd.RootCmd, "./")
}
