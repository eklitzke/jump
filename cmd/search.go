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

	"github.com/eklitzke/jump/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the database for matches",
	Run: func(cmd *cobra.Command, args []string) {
		handle := db.LoadDefaultDatabase()
		entry := handle.Search(args[0])
		if entry.Path != "" {
			log.Debug().Str("path", entry.Path).Float64("weight", entry.Weight).Msg("found match")
			fmt.Println(entry.Path)
			return
		}
		log.Debug().Msg("no match found")
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
