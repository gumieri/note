package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
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

func getTextFromEditor(editorCommand string, fileName string) (text string, err error) {
	filePath := filepath.Join(os.TempDir(), fileName)

	tmpFile, err := os.Create(filePath)
	if err != nil {
		return
	}

	tmpFile.Close()

	editorCmd := exec.Command(editorCommand, filePath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	err = editorCmd.Start()
	if err != nil {
		return
	}

	err = editorCmd.Wait()
	if err != nil {
		return
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	text = string(content)

	return
}

func writeNote(context *cli.Context) {
	notePath := os.Getenv("NOTE_PATH")

	nextNumber, err := incrementNoteNumber(notePath)

	if err != nil {
		log.Fatal(err)
	}

	nextNumberString := strconv.Itoa(nextNumber)

	var noteContent string
	if len(context.Args()) == 0 {
		noteContent, err = getTextFromEditor(os.Getenv("EDITOR"), nextNumberString)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		noteContent = strings.Join(context.Args()[:], " ") + "\n"
	}

	noteTitle := noteContent

	if noteContent == "" {
		log.Fatal(errors.New("Empty content"))
	}

	firstLineBreakIndex := strings.Index(noteTitle, "\n")
	if firstLineBreakIndex >= 0 {
		noteTitle = noteTitle[0:firstLineBreakIndex]
	}

	if len(noteContent) > 72 {
		noteTitle = noteContent[0:72]

		nextCharacter := noteContent[72:73]
		if nextCharacter != " " && nextCharacter != "\n" {
			lastSpace := strings.LastIndex(noteTitle, " ")
			noteTitle = noteTitle[0:lastSpace]
		}
	}

	noteName := nextNumberString + " - " + noteTitle

	_ = os.Mkdir(notePath, 0755)

	noteFile, err := os.Create(filepath.Join(notePath, noteName))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(noteName)
	noteFile.WriteString(noteContent)

	defer noteFile.Close()

	return
}

func showNote(context *cli.Context) {
	notePath := os.Getenv("NOTE_PATH")

	notesNames, err := existingNotesNames(notePath)

	if err != nil {
		log.Fatal(err)
	}

	noteToFind := strings.Join(context.Args()[:], " ")

	notesFound := fuzzy.RankFind(noteToFind, notesNames)
	noteFound := notesFound[0].Target

	noteContent, err := ioutil.ReadFile(filepath.Join(notePath, noteFound))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(noteFound + "\n")
	fmt.Print(string(noteContent))
	return
}

func deleteNote(context *cli.Context) {
	notePath := os.Getenv("NOTE_PATH")

	notesNames, err := existingNotesNames(notePath)

	if err != nil {
		log.Fatal(err)
	}

	noteToFind := strings.Join(context.Args()[:], " ")

	notesFound := fuzzy.RankFind(noteToFind, notesNames)
	noteFound := notesFound[0].Target

	err = os.Remove(filepath.Join(notePath, noteFound))

	if err != nil {
		log.Fatal(err)
	}

	return
}

func listNotes(context *cli.Context) {
	notePath := os.Getenv("NOTE_PATH")

	notesNames, err := existingNotesNames(notePath)

	if err != nil {
		log.Fatal(err)
	}

	for _, note := range notesNames {
		fmt.Println(note)
	}

	return
}

func main() {
	app := cli.NewApp()

	app.Name = "Note"

	app.Version = "0.0.1"

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
			Name:    "delete",
			Aliases: []string{"del", "rm"},
			Usage:   "delete a note",
			Action:  deleteNote,
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list notes",
			Action:  listNotes,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
