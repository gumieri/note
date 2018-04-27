package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gumieri/note/lib/notes"
	editor "github.com/gumieri/open-in-editor"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// EditNote accept one argument to execute a
// fuzzy search for the note to be edited
// and open the EDITOR with the found result
func EditNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	noteFound, err := notes.FindNoteName(notePath, context.Args()[:])

	err = editor.File(viper.GetString("editor"), filepath.Join(notePath, noteFound))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}
