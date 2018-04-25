package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	editor "github.com/gumieri/open-in-editor"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func existingNotesNames(notePath string) (notesNames []string, err error) {
	f, err := os.Open(notePath)

	if err != nil {
		return
	}

	list, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		return
	}

	sort.Slice(list, func(a, b int) bool { return list[a].Name() < list[b].Name() })

	for _, file := range list {
		notesNames = append(notesNames, file.Name())
	}

	return
}

func incrementNoteNumber(notePath string) (number int, err error) {
	notesNames, err := existingNotesNames(notePath)

	if err != nil {
		return
	}

	if len(notesNames) == 0 {
		return
	}

	lastNoteName := notesNames[len(notesNames)-1]
	re := regexp.MustCompile("^[0-9]+")
	number, err = strconv.Atoi(re.FindAllString(lastNoteName, 1)[0])

	if err != nil {
		return
	}

	number = number + 1

	return
}

func writeNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	nextNumber, err := incrementNoteNumber(notePath)

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

func showNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	notesNames, err := existingNotesNames(notePath)

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

	fmt.Println(noteFound + "\n")
	fmt.Print(string(noteContent))
	return
}

func deleteNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	notesNames, err := existingNotesNames(notePath)

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

func editNote(context *cli.Context) {
	notePath := viper.GetString("notePath")

	notesNames, err := existingNotesNames(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	noteToFind := strings.Join(context.Args()[:], " ")

	notesFound := fuzzy.RankFind(noteToFind, notesNames)
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

func listNotes(context *cli.Context) {
	notePath := viper.GetString("notePath")

	notesNames, err := existingNotesNames(notePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, note := range notesNames {
		fmt.Println(note)
	}

	return
}

func main() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetDefault("editor", "vim")
	viper.SetDefault("notePath", filepath.Join(currentUser.HomeDir, "Notes"))

	viper.SetConfigName(".noteconfig")
	viper.AddConfigPath(currentUser.HomeDir)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	app := cli.NewApp()

	app.Name = "Note"

	app.Version = "0.0.3"

	app.Usage = "Quick and easy Command-line tool for taking notes"
	app.UsageText = "note [just type a text] [or command] [with command options]"
	app.ArgsUsage = "[text]"

	app.Action = writeNote

	app.Commands = []cli.Command{
		{
			Name:   "show",
			Usage:  "show a note contet",
			Action: showNote,
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "edit a note contet",
			Action:  editNote,
		},
		{
			Name:    "delete",
			Aliases: []string{"del", "d", "rm"},
			Usage:   "delete a note",
			Action:  deleteNote,
		},
		{
			Name:    "list",
			Aliases: []string{"ls", "l"},
			Usage:   "list notes",
			Action:  listNotes,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
