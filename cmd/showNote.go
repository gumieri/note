package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/gumieri/note/lib/notes"
)

// ShowNote show the content of a note
// it will do a fuzzy search using the argument for the note name
func ShowNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	notesNames, err := notes.ExistingNames(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	noteToFind := strings.Join(context.Args()[:], " ")

	notesFound := fuzzy.RankFind(noteToFind, notesNames)
	noteFound := notesFound[0].Target

	noteContent, err := ioutil.ReadFile(filepath.Join(notePath, noteFound))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s\n\n%s", noteFound, noteContent)
	return
}
