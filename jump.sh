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

# A simple jumper.
jump_jump() {
  local dest
  dest="$(jump search "$1")"
  if [[ -n "$dest" ]] && [[ -d "$dest" ]]; then
    printf '\033[0;31m%s\033[0m\n' "$dest"
    cd "$dest"
  fi
}

# If the jump command is available, alias j=jump_jump and add "jump update" to
# the PROMPT_COMMAND.
if command -v jump &>/dev/null; then
  if [[ "x${_JUMP_ENABLED}" = x ]]; then
    PROMPT_COMMAND="${PROMPT_COMMAND};jump update"
    _JUMP_ENABLED=yes
  fi
  alias j=jump_jump
fi
