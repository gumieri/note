package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	editor "github.com/gumieri/open-in-editor"
	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/gumieri/note/lib/notes"
)

// WriteNote when no argument is given
// it open the configured EDITOR
// When arguments are informed
// these became the contet of the note
func WriteNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	_, err := os.Stat(notePath)

	if err != nil {
		err := os.Mkdir(notePath, 0755)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	noteTitle := context.String("title")

	nextNumber, err := notes.NextNumber(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var noteContent string
	if len(context.Args()) == 0 {
		editorCommand := viper.GetString("editor")

		var tmpTitle string
		if noteTitle == "" {
			tmpTitle = "new note"
		} else {
			tmpTitle = noteTitle
		}

		tmpFileName := notes.NoteName(nextNumber, tmpTitle)
		noteContent, err = editor.GetContentFromTemporaryFile(editorCommand, tmpFileName, "")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		noteContent = fmt.Sprintf("%s\n", strings.Join(context.Args(), " "))
	}

	if noteContent == "" {
		fmt.Println(errors.New("Empty content"))
		os.Exit(1)
	}

	if noteTitle == "" {
		noteTitle = notes.FormatTitle(noteContent)
	}

	noteName := notes.NoteName(nextNumber, noteTitle)

	noteFile, err := os.Create(filepath.Join(notePath, noteName))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	noteFile.WriteString(noteContent)

	defer noteFile.Close()

	fmt.Printf("%d\t%s\n", nextNumber, noteTitle)

	return
}
