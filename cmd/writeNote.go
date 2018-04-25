package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	nextNumber, err := notes.NextNumber(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nextNumberString := strconv.Itoa(nextNumber)

	editorCommand := viper.GetString("editor")
	tmpFileName := nextNumberString + " - new note"
	var noteContent string
	if len(context.Args()) == 0 {
		noteContent, err = editor.GetContentFromTemporaryFile(editorCommand, tmpFileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		noteContent = strings.Join(context.Args()[:], " ") + "\n"
	}

	if noteContent == "" {
		fmt.Println(errors.New("Empty content"))
		os.Exit(1)
	}

	noteTitle := strings.Split(noteContent, "\n")[0]
	if len(noteTitle) > 72 {
		nextCharacter := noteTitle[72:73]

		noteTitle = noteTitle[0:72]

		if nextCharacter != " " {
			lastSpace := strings.LastIndex(noteTitle, " ")
			noteTitle = noteTitle[0:lastSpace]
		}
	}

	noteName := nextNumberString + " - " + noteTitle

	_ = os.Mkdir(notePath, 0755)

	noteFile, err := os.Create(filepath.Join(notePath, noteName))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(noteName)
	noteFile.WriteString(noteContent)

	defer noteFile.Close()

	return
}
