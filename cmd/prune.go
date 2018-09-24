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
	"github.com/spf13/cobra"
)

var pruneCount int

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Automatically prune old or invalid database entries",
	Run: func(cmd *cobra.Command, args []string) {
		handle.Prune(pruneCount)
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
	pruneCmd.Flags().IntVarP(&pruneCount, "num-database-entries", "n", 1000, "Number of database entries to keep")
}
