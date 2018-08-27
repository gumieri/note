package notes

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

var noteRegex = regexp.MustCompile("^[0-9]+ - ")
var numberRegex = regexp.MustCompile("^[0-9]+")

// TitleFromNoteName extract title from a noteName
func TitleFromNoteName(noteName string) string {
	return noteRegex.ReplaceAllString(noteName, "")
}

// NumberFromNoteName extract number from a noteName
func NumberFromNoteName(noteName string) int {
	number, _ := strconv.Atoi(numberRegex.FindAllString(noteName, 1)[0])
	return number
}

// NoteName join number and title
func NoteName(number int, title string) string {
	return fmt.Sprintf("%d - %s", number, title)
}

// FormatTitle return the first line of a string and
// cut the words to be smaller than 72 characters
func FormatTitle(raw string) (formated string) {
	formated = strings.Split(raw, "\n")[0]

	if len(formated) <= 72 {
		return
	}

	nextCharacter := formated[72:73]

	formated = formated[0:72]

	if nextCharacter == " " {
		return
	}

	lastSpace := strings.LastIndex(formated, " ")
	formated = formated[0:lastSpace]
	return
}

// FindNoteName return a note from slice of string
func FindNoteName(notePath string, words []string, caseSensitive bool) (noteName string, err error) {
	notesNames, err := ExistingNames(notePath)

	if err != nil {
		return
	}

	var notesFound fuzzy.Ranks
	if caseSensitive {
		notesFound = fuzzy.RankFind(strings.Join(words, " "), notesNames)
	} else {
		notesFound = fuzzy.RankFindFold(strings.Join(words, " "), notesNames)
	}

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

	sort.Slice(list, func(a, b int) bool {
		return NumberFromNoteName(list[a].Name()) < NumberFromNoteName(list[b].Name())
	})

	for _, file := range list {
		if noteRegex.MatchString(file.Name()) {
			notesNames = append(notesNames, file.Name())
		}
	}

	return
}

// NextNumber return the incremented number from ExistingNames
func NextNumber(notePath string) (number int, err error) {
	notesNames, err := ExistingNames(notePath)

	if err != nil || len(notesNames) == 0 {
		return
	}

	lastNoteName := notesNames[len(notesNames)-1]

	number = NumberFromNoteName(lastNoteName) + 1

	return
}
