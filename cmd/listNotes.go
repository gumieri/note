package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/gumieri/note/lib/notes"
)

// ListNotes list all existing notes in the NOTE_PATH
func ListNotes(context *cli.Context) {
	notePath := viper.GetString("notePath")

	notesNames, err := notes.ExistingNames(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, note := range notesNames {
		fmt.Println(note)
	}

	return
}
