package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	editor "github.com/gumieri/open-in-editor"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/gumieri/note/lib/notes"
)

// EditNote accept one argument to execute a
// fuzzy search for the note to be edited
// and open the EDITOR with the found result
func EditNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	notesNames, err := notes.ExistingNames(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	noteToFind := strings.Join(context.Args()[:], " ")

	notesFound := fuzzy.RankFind(noteToFind, notesNames)

	if len(notesFound) == 0 {
		os.Exit(1)
	}

	noteFound := notesFound[0].Target

	err = editor.File(viper.GetString("editor"), filepath.Join(notePath, noteFound))
	if err != nil {
		return
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}
