// Copyright 2018 Evan Klitzke <evan@eklitzke.org>
//
// This file is part of jump.
//
// jump is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// jump is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// jump. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"github.com/eklitzke/jump/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var dumpShort bool

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump database contents as plaintext",
	Run: func(cmd *cobra.Command, args []string) {
		dumpOpts := db.DumpOpts{
			Short: dumpShort,
		}
		obj := db.Dump(handle, dumpOpts)
		enc := newStdoutJSONEncoder()
		if err := enc.Encode(obj); err != nil {
			log.Warn().Err(err).Msg("failed to json encode database")
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().BoolVar(&dumpShort, "short-paths", false, "Dump shortened paths")
}
