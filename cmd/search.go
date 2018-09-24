// Copyright 2019 Evan Klitzke <evan@eklitzke.org>
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
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var verbose bool
var searchCount int

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the database for matches",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal().Int("argcount", len(args)).Interface("args", args).Msg("expected exactly one search argument")
		}
		entries := handle.Search(searchCount, args...)
		for _, entry := range entries {
			if verbose {
				fmt.Printf("%10.4f  %s\n", entry.Weight, entry.Path)
			} else {
				fmt.Println(entry.Path)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&searchCount, "num-results", "n", 1, "Number of database entries to keep")
	searchCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose results")
}
