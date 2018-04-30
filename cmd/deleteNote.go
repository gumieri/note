package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	libCli "github.com/gumieri/note/lib/cli"
	"github.com/gumieri/note/lib/notes"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// DeleteNote accept an argument to execute
// a fuzzy search for the note to be deleted
func DeleteNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	_, err := os.Stat(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	noteFound, err := notes.FindNoteName(notePath, context.Args()[:])

	confirmMessage := fmt.Sprintf("Do you want to delete the note %s ?", noteFound)
	if !context.Bool("yes") && !libCli.Confirm(confirmMessage) {
		return
	}

	err = os.Remove(filepath.Join(notePath, noteFound))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}
