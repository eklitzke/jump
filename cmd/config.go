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
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	ExcludePatterns []string `yaml:"ExcludePatterns"`
}

func loadConfig() *config {
	c := &config{}
	yamlFile, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return c
	}
	_ = yaml.Unmarshal(yamlFile, c)
	return c
}
