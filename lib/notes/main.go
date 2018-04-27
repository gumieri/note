package notes

import (
	"errors"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

// FindNoteName return a note from slice of string
func FindNoteName(notePath string, words []string) (noteName string, err error) {
	notesNames, err := ExistingNames(notePath)

	if err != nil {
		return
	}

	notesFound := fuzzy.RankFind(strings.Join(words, " "), notesNames)

	if len(notesFound) == 0 {
		err = errors.New("No note found")
		return
	}

	noteName = notesFound[0].Target
	return
}

// ExistingNames return the file names from the given NOTE_PATH
func ExistingNames(notePath string) (notesNames []string, err error) {
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

// NextNumber return the incremented number from ExistingNames
func NextNumber(notePath string) (number int, err error) {
	notesNames, err := ExistingNames(notePath)

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
