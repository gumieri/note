# Note

Quick and easy Command-line tool for taking notes

## Installation

1. Download the executable for your Operating System at [note/releases](https://github.com/gumieri/note/releases/latest);
2. Rename the executable to `note`;
3. Place the executable in a directory loaded by the system ([about these directories](https://en.wikipedia.org/wiki/PATH_%28variable%29));
4. Give it permission to be executable (only if you are using Mac OS or Linux).

Or, if you are using a Mac OS or Linux, you can just execute the commands described at the [note/releases](https://github.com/gumieri/note/releases/latest).

## Configurations

The default configuration is:
```yml
editor: vim
notePath: ~/Notes
```

You can create a `.noteconfig.yml` (or json, or toml) in your home directory to o override these configurations.
If you create a configuration file in a specific directory, it will take priority over the default and the configuration in the home directory.
As well the `EDITOR` and `NOTE_PATH` environment variables has priority over these configuration files.

## Usage

```bash
note [just type a text] [or command] [with command options]
```

Note is very easy and simple to use.
Start by typing `note` command and continue describing what you want to take as note:
```bash
note there is no place like home
```

If using special character, just use quotes or escape the character:
```bash
note "there's no place like home"
```
```bash
note there\'s no place like home
```

It will create a note with the title `0 - there's no place like home`.
In case of you need a text editor you can just type `note` without any argument.
It will open the `EDITOR` defined as environment variable or the configured one.

#### Title

To define a title, just use the flag option `--title` (or `-t`).
If no title is informed, `note` will take some words from the first line as it.
All titles start with a number (integer) increasing by one from the last note for better identification.

### Other commands

#### Show

Show a note content.
It will search for a note using the given arguments executing a fuzzy search:
```
note show like home
```
If it's your first note, you can surely show it by typing:
```
note show 0
```

#### Edit

Edit a note content.
Like the show command, will use the given arguments to search for a note but will open it content in your text editor.
```
note edit 0
```
```
note e 0
```

#### Delete

Delete a note.
Like show and edit, but delete a note.
It will ask for confirmation if not given the flag option `--yes` / `-y`.
```
note delete 0
```
```
note del 0
```
```
note d 0
```
```
note rm 0
```

#### List
List notes.
Has no arguments. Simply list the notes at the `notePath`.
```
note list
```
```
note ls
```
```
note l
```

## Code Status

[![Go Report Card](https://goreportcard.com/badge/github.com/gumieri/note)](https://goreportcard.com/report/github.com/gumieri/note)
[![Build Status](https://travis-ci.org/gumieri/note.svg?branch=master)](https://travis-ci.org/gumieri/note)


## License

Note is released under the [MIT License](http://www.opensource.org/licenses/MIT).

