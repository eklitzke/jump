#!/bin/bash
#
# Copyright 2019 Evan Klitzke <evan@eklitzke.org>
#
# This file is part of jump.
#
# jump is free software: you can redistribute it and/or modify it under
# the terms of the GNU General Public License as published by the Free Software
# Foundation, either version 3 of the License, or (at your option) any later
# version.
#
# jump is distributed in the hope that it will be useful, but WITHOUT ANY
# WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
# A PARTICULAR PURPOSE. See the GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License along with
# jump. If not, see <http://www.gnu.org/licenses/>.

_print_red() {
  printf '\033[0;31m%s\033[0m\n' "$1"
}

# Try to jump to the best matching entry in the jump database.
j() {
  # Print help if that's the search query (use "j -- help" to use "help" as the
  # actual query).
  if [[ $# -eq 1 ]] && [[ "$1" == help ]]; then
    echo "Usage:"
    echo "  j QUERY     jump to directory matching QUERY"
    echo "  jc QUERY    jump to subdirectory matching QUERY"
    echo "  jo QUERY    open the file matching QUERY"
    echo "  jco QUERY   open the subdirectory file matching QUERY"
    return
  fi

  local dest
  dest="$(jump search "$@")"
  if [[ -n "$dest" ]] ; then
    _print_red "$dest"
    cd "$dest" || return 1
  else
    _print_red "no matches found"
  fi
}

# Jump to child directory.
jc() {
  j "$PWD" "$@"
}

# Open a file using xdg-open. This is kind of stupid, but it's provided here to
# provide feature parity with autojump.
jo() {
  local f
  f="$(jump search "$@")"
  if [[ -f "$f" ]]; then
    _print_red "$f"
    xdg-open "$f"
  else
    _print_red "no matches found"
  fi
}

# Likewise, but for the child directory.
jco() {
  jc "$PWD" "$@"
}

# Check if jump is available, and if so set up PROMPT_COMMAND.
if command -v jump &>/dev/null; then
  if [[ "x${JUMP_ENABLED}" = x ]]; then
    PROMPT_COMMAND="${PROMPT_COMMAND};jump update"
    JUMP_ENABLED=yes
    export JUMP_ENABLED
  fi
fi
