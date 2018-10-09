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
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var varsCmd = &cobra.Command{
	Use:   "vars",
	Short: "Print variables",
	Run: func(cmd *cobra.Command, args []string) {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(struct {
			Config   string `json:"config"`
			Database string `json:"database"`
		}{
			Config:   cfgFile,
			Database: dbPath,
		}); err != nil {
			log.Warn().Err(err).Msg("failed to json encode vars")
		}
	},
}

func init() {
	rootCmd.AddCommand(varsCmd)
}
