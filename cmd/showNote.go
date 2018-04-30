package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gumieri/note/lib/notes"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// ShowNote show the content of a note
// it will do a fuzzy search using the argument for the note name
func ShowNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	_, err := os.Stat(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	noteFound, err := notes.FindNoteName(notePath, context.Args()[:])

	noteContent, err := ioutil.ReadFile(filepath.Join(notePath, noteFound))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s\n\n%s", noteFound, noteContent)
	return
}
