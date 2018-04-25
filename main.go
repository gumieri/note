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

	app.Version = "0.0.4"

	app.Usage = "Quick and easy Command-line tool for taking notes"
	app.UsageText = "note [just type a text] [or command] [with command options]"
	app.ArgsUsage = "[text]"

	app.Action = cmd.WriteNote

	app.Commands = []cli.Command{
		{
			Name:   "show",
			Usage:  "show a note contet",
			Action: cmd.ShowNote,
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "edit a note contet",
			Action:  cmd.EditNote,
		},
		{
			Name:    "delete",
			Aliases: []string{"del", "d", "rm"},
			Usage:   "delete a note",
			Action:  cmd.DeleteNote,
		},
		{
			Name:    "list",
			Aliases: []string{"ls", "l"},
			Usage:   "list notes",
			Action:  cmd.ListNotes,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
