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

func editContent(notePath string, words []string) (err error) {
	noteFound, err := notes.FindNoteName(notePath, words)

	if err != nil {
		return
	}

	err = editor.File(viper.GetString("editor"), filepath.Join(notePath, noteFound))

	return
}

func editTitle(notePath string, newTitle string, words []string) (err error) {
	if len(words) == 0 {
		words = append(words, newTitle)
		newTitle = ""
	}

	oldName, err := notes.FindNoteName(notePath, words)

	if err != nil {
		return
	}

	if newTitle == "" {
		oldTitle := notes.TitleFromNoteName(oldName)

		newTitle, err = editor.GetContentFromTemporaryFile(viper.GetString("editor"), oldName, oldTitle)

		if err != nil {
			return
		}
	}

	newName := notes.NoteName(notes.NumberFromNoteName(oldName), notes.FormatTitle(newTitle))
	err = os.Rename(filepath.Join(notePath, oldName), filepath.Join(notePath, newName))

	return
}

// EditNote accept one argument to execute a
// fuzzy search for the note to be edited
// and open the EDITOR with the found result
func EditNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	noteTitle := context.String("title")

	_, err := os.Stat(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if noteTitle == "" {
		err = editContent(notePath, context.Args())
	} else {
		err = editTitle(notePath, noteTitle, context.Args())
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
