[![CircleCI](https://circleci.com/gh/eklitzke/jump/tree/master.svg?style=shield)](https://circleci.com/gh/eklitzke/jump/tree/master)
[![Documentation](https://godoc.org/github.com/eklitzke/jump/db?status.svg)](http://godoc.org/github.com/eklitzke/jump/db)
[![codecov](https://codecov.io/gh/eklitzke/jump/branch/master/graph/badge.svg)](https://codecov.io/gh/eklitzke/jump)

This project implements a simple shell command jumper similar to
[autojump](https://github.com/wting/autojump) or
[fasd](https://github.com/clvv/fasd). My version is just called `jump`, because
the good names were already taken. I wrote this for my own personal use, to work
around defects I found in both of those projects. This implementation is
probably faster than the alternatives (it's written in Go and uses a binary
database format), but it has fewer features and I doubt the performance
difference is noticeable anyway. I think `jump` is less likely to silently
corrupt your jump database than autojump, but as always the principle of *caveat
emptor* applies.

Features:

 * Small and fast, written in Go
 * Binary database format (proven to be 666% faster than text databases)
 * Incorporates access recency in search rankings
 * Compatibility with autojump shell commands (`j`, `jc`, `jo`, and `jco`)

## Installation

If you use Bash and want to try out my code, you can set it up like this:

```bash
# Install the jump command
$ go get -u github.com/eklitzke/jump

# Source jump.sh in your .bashrc
$ curl -sL https://raw.githubusercontent.com/eklitzke/jump/master/jump.sh -o ~/.jump.sh
$ echo '. ~/.jump.sh' >> ~/.bashrc
```

To check that everything is set up correctly, launch a new Bash shell (e.g. by
creating a new terminal window) and check that you see output something like
this:

```bash
# Check that j is available as a shell function
$ type -t j
function
```

## Usage

Use the `j` command to jump places. For example, if you run a lot of commands in
a directory named `~/foo/bar`, running the shell command `j bar` should jump to
the `~/foo/bar` directory.

For more advanced commands run `jump help`:

```plain
$ jump help
Jump is a shell autojumper

Usage:
  jump [command]

Available Commands:
  dump        Dump database contents as plaintext
  help        Help about any command
  import      Import an autojump database
  prune       Automatically prune old or invalid database entries
  remove      Remove a database entry
  search      Search the database for matches
  update      Update database weights

Flags:
  -c, --config string      config file (default "/home/evan/.config/jump/jump.yml")
  -D, --database string    database file (default "/home/evan/.local/share/jump/db.gob")
  -d, --debug              enable debug mode
  -h, --help               help for jump
      --log-caller         include caller info in log messages
  -l, --log-level string   the log level (default "info")
      --time-matching      enable time matching in searches (default true)

Use "jump [command] --help" for more information about a command.
```

### Issues With `PROMPT_COMMAND`

The `jump.sh` shell code makes use of `PROMPT_COMMAND` in order to maintain the
jump database. That means that blindly overwriting `PROMPT_COMMAND` elsewhere in
your Bash profile will cause `jump` to stop working. If you want to set your own
`PROMPT_COMMAND` all you need to do is make sure you append to the variable
rather than overwriting it (you can look at `jump.sh` itself for an example of
how to do this correctly).
