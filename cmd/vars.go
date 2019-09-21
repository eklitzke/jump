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

func newStdoutJSONEncoder() *json.Encoder {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc
}

type configDisplay struct {
	Paths           map[string]string `json:"paths"`
	ExcludePatterns []string          `json:"excludePatterns"`
}

var varsCmd = &cobra.Command{
	Use:   "vars",
	Short: "Print variables",
	Run: func(cmd *cobra.Command, args []string) {
		c := loadConfig()
		display := configDisplay{
			Paths: map[string]string{
				"config":   cfgFile,
				"database": dbPath,
			},
			ExcludePatterns: c.ExcludePatterns,
		}
		enc := newStdoutJSONEncoder()
		if err := enc.Encode(display); err != nil {
			log.Warn().Err(err).Msg("failed to json encode vars")
		}
	},
}

func init() {
	rootCmd.AddCommand(varsCmd)
}
