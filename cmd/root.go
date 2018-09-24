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
	"os"

	"github.com/eklitzke/jump/db"
	isatty "github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var dbPath string
var logCaller bool
var logLevel string
var handle *db.Database

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jump",
	Short: "Jump is a shell autojumper",
	// TODO: make the raw command equivalent to search?
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("fatal error running command")
	}

	// Save the database; this is a no-op if the database hasn't been
	// mutated.
	if err := handle.Save(); err != nil {
		log.Fatal().Err(err).Msg("failed to save database")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initDBHandle)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", db.ConfigPath(), "config file")
	rootCmd.PersistentFlags().StringVarP(&dbPath, "database", "d", db.DatabasePath(), "database file")
	rootCmd.PersistentFlags().BoolVar(&logCaller, "log-caller", false, "include caller info in log messages")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "the log level")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Str("file", viper.ConfigFileUsed()).Msg("read config file")
	}
}

func initLogging() {
	if isatty.IsTerminal(os.Stderr.Fd()) {
		w := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = zerolog.New(w).With().Timestamp().Logger()
	}
	if logCaller {
		log.Logger = log.Logger.With().Caller().Logger()
	}
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Debug().Str("level", level.String()).Msg("logging initialized")

}

func initDBHandle() {
	handle = db.LoadDatabase(dbPath)
}
