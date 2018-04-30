package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gumieri/note/lib/notes"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// ListNotes list all existing notes in the NOTE_PATH
func ListNotes(context *cli.Context) {
	notePath := viper.GetString("notePath")

	_, err := os.Stat(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	notesNames, err := notes.ExistingNames(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	color.New(color.Bold).Println("ID\tTitle")
	for _, noteName := range notesNames {
		number := notes.NumberFromNoteName(noteName)
		title := notes.TitleFromNoteName(noteName)
		fmt.Printf("%d\t%s\n", number, title)
	}

	return
}
