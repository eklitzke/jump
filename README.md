[![CircleCI](https://circleci.com/gh/eklitzke/jump/tree/master.svg?style=svg)](https://circleci.com/gh/eklitzke/jump/tree/master)

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

## Installation

If you use Bash and want to try out my code, you can set it up like this:

```bash
# Install the jump command
$ go get github.com/eklitzke/jump

# Source jump.sh in your .bashrc
$ curl -sL https://raw.githubusercontent.com/eklitzke/jump/master/jump.sh -o ~/.jump.sh
$ echo '. ~/.jump.sh' >> ~/.bashrc
```

To check that everything is set up correctly, launch a new Bash shell (e.g. by
creating a new terminal window) and check that you see output something like
this:

```bash
# Check that j is aliased to the jump.sh provided jump_jump function.
$ type j
j is aliased to `jump_jump'
```

## Usage

Use the `j` command to jump places. For example, if you run a lot of commands in
a directory named `~/foo/bar`, running the shell command `j bar` should jump to
the `~/foo/bar` directory.

For more advanced commands run `jump help`.

### Issues With `PROMPT_COMMAND`

The `jump.sh` shell code makes use of `PROMPT_COMMAND` in order to maintain the
jump database. That means that blindly overwriting `PROMPT_COMMAND` elsewhere in
your Bash profile will cause `jump` to stop working. If you want to set your own
`PROMPT_COMMAND` all you need to do is make sure you append to the variable
rather than overwriting it (you can look at `jump.sh` itself for an example of
how to do this correctly).
