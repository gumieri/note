package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/gumieri/note/lib/notes"
)

// DeleteNote accept an argument to execute
// a fuzzy search for the note to be deleted
func DeleteNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	notesNames, err := notes.ExistingNames(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	noteToFind := strings.Join(context.Args()[:], " ")

	notesFound := fuzzy.RankFind(noteToFind, notesNames)
	noteFound := notesFound[0].Target

	err = os.Remove(filepath.Join(notePath, noteFound))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}