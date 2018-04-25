package notes

import (
	"os"
	"regexp"
	"sort"
	"strconv"
)

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
