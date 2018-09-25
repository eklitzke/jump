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
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/eklitzke/jump/db"
	isatty "github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var dbPath string
var debug bool
var timeMatching bool
var logCaller bool
var logLevel string
var handle db.Database

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
	if handle != nil {
		if err := saveDB(); err != nil {
			log.Fatal().Err(err).Msg("failed to save database")
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initDBHandle)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", db.ConfigPath(), "config file")
	rootCmd.PersistentFlags().StringVarP(&dbPath, "database", "D", db.DatabasePath(), "database file")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	rootCmd.PersistentFlags().BoolVar(&logCaller, "log-caller", false, "include caller info in log messages")
	rootCmd.PersistentFlags().BoolVar(&timeMatching, "time-matching", true, "enable time matching in searches")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "the log level")

	// Start logging initialization now, so that log messages are properly
	// formatted on the console if other initialization tasks fail.
	if isatty.IsTerminal(os.Stderr.Fd()) {
		w := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = zerolog.New(w).With().Timestamp().Logger()
	}
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
	if logCaller {
		log.Logger = log.Logger.With().Caller().Logger()
	}
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	if debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Debug().Str("level", level.String()).Msg("logging initialized")

}

func initDBHandle() {
	var r io.Reader
	dbFile, err := os.Open(dbPath)
	if err != nil {
		if !os.IsNotExist(err) {
			// a serious error
			log.Fatal().Err(err).Str("path", dbPath).Msg("failed to open database file")
		}
		// not so serious
		log.Debug().Err(err).Str("path", dbPath).Msg("database file not found")
		r = &bytes.Buffer{}
	} else {
		defer func() {
			if err := dbFile.Close(); err != nil {
				log.Warn().Err(err).Str("path", dbPath).Msg("failed to close db file")
			}
		}()
		r = dbFile
	}
	handle = db.NewDatabase(r, db.Options{
		Debug:        debug,
		TimeMatching: timeMatching,
	})
}

func saveDB() error {
	if !handle.Dirty() {
		log.Debug().Msg("database not dirty, skipping save")
		return nil
	}

	// ensure the directory exists
	dir := filepath.Dir(dbPath)
	ensureDirectory(dir)

	// create the temporary file in the same directory as the destination
	// file, to ensure that the rename operation is atomic
	temp, err := ioutil.TempFile(dir, ".jump.bak")
	if err != nil {
		log.Error().Err(err).Str("dir", dir).Msg("failed to create temporary save file")
		return err
	}

	// clean up the temporary file when we're done with it
	tempName := temp.Name()
	defer func() {
		if err := temp.Close(); err != nil {
			log.Error().Err(err).Str("path", tempName).Msg("error closing temporary file")
		}
		if err := os.Remove(tempName); err != nil && !os.IsNotExist(err) {
			log.Error().Err(err).Str("path", tempName).Msg("failed to close temporary file")
		}
	}()

	// encode and flush the file
	w := bufio.NewWriter(temp)
	if err := handle.Save(w); err != nil {
		log.Error().Err(err).Msg("failed to encode database")
		return err
	}
	if err := w.Flush(); err != nil {
		log.Error().Err(err).Msg("failed to flush temporary file")
		return err
	}

	// atomic rename
	if err := os.Rename(tempName, dbPath); err != nil {
		log.Error().Err(err).Str("dbpath", dbPath).Str("tempfile", tempName).Msg("failed to rename db file")
		return err
	}

	return nil
}

// ensureDirectory ensures that a directory exists
func ensureDirectory(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// create the data directory
			if err := os.Mkdir(dir, 0700); err != nil {
				log.Fatal().Err(err).Str("dir", dir).Msg("failed to create directory")
			}
		} else {
			log.Fatal().Err(err).Str("dir", dir).Msg("failed to stat directory")
		}
	}
}
