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
  if [[ "$1" == -* ]] && [[ "$1" != "--" ]]; then
    j "$@"
  else
    j "$PWD" "$@"
  fi
}

# N.B. jo and jco are not defined because they shouldn't exist in the first
# place.

# Check if jump is available, and if so set up PROMPT_COMMAND.
if command -v jump &>/dev/null; then
  if [[ "x${JUMP_ENABLED}" = x ]]; then
    PROMPT_COMMAND="${PROMPT_COMMAND};jump update"
    JUMP_ENABLED=yes
    export JUMP_ENABLED
  fi
fi
