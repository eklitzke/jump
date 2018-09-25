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
	"fmt"

	"github.com/eklitzke/jump/db"
	"github.com/spf13/cobra"
)

var verbose bool
var searchCount int

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the database for matches",
	Run: func(cmd *cobra.Command, args []string) {
		var printer func(db.Entry)
		if verbose {
			printer = func(e db.Entry) { fmt.Printf("%10.4f  %s\n", e.Weight, e.Path) }
		} else {
			printer = func(e db.Entry) { fmt.Println(e.Path) }
		}
		entries := handle.Search(searchCount, args...)
		for _, entry := range entries {
			printer(entry)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&searchCount, "num-results", "n", 1, "Number of database entries to keep")
	searchCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose results")
}
