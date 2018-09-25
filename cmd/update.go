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
	"os"

	"github.com/eklitzke/jump/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var updateWeight float64

// updateCmd represents the add command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update database weights",
	Run: func(cmd *cobra.Command, args []string) {
		if updateWeight == 0 {
			log.Fatal().Msg("ignoring update command for 0 weight")
		}
		if len(args) == 0 {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal().Err(err).Msg("no argument supplied and failed to getcwd")
			}
			args = append(args, dir)
		}

		// try to update each argument, first checking that it exists and is a directory
		for _, dir := range args {
			if err := db.CheckIsDir(dir); err != nil {
				continue
			}
			handle.AdjustWeight(dir, updateWeight)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Float64VarP(&updateWeight, "weight", "w", 15, "Weight to adjust by (may be negative)")
}
