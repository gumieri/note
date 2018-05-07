package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/gumieri/note/cmd"
)

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

	app.Version = "1.1.0"

	app.Usage = "Quick and easy Command-line tool for taking notes"
	app.UsageText = "note [just type a text] [or command] [with command options]"
	app.ArgsUsage = "[text]"

	app.Action = cmd.WriteNote

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "title, t",
			Usage: "Inform a title for the note",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "show",
			Usage:  "Show a note content",
			Action: cmd.ShowNote,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-title",
					Usage: "Don't print note's title",
				},
				cli.BoolFlag{
					Name:  "case-sensitive, s",
					Usage: "Must match the case of the informed parameters and note's title",
				},
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "Edit a note content",
			Action:  cmd.EditNote,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "title, t",
					Usage: "Edit the title of the note",
				},
				cli.BoolFlag{
					Name:  "case-sensitive, s",
					Usage: "Must match the case of the informed parameters and note's title",
				},
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"del", "d", "rm"},
			Usage:   "Delete a note",
			Action:  cmd.DeleteNote,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yes, y",
					Usage: "Delete without asking confirmation",
				},
				cli.BoolFlag{
					Name:  "case-sensitive, s",
					Usage: "Must match the case of the informed parameters and note's title",
				},
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls", "l"},
			Usage:   "List notes",
			Action:  cmd.ListNotes,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-header",
					Usage: "Don't print the description header of the columns",
				},
				cli.BoolFlag{
					Name:  "filename",
					Usage: "Print only the filenames of the notes at the notePath",
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
