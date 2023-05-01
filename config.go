/*
 *  Copyright 2023, Enguerrand de Rochefort
 *
 * This file is part of left.
 *
 * left is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * left is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with left.  If not, see <http://www.gnu.org/licenses/>.
 *
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

type Config struct {
	FontName          string
	FontSize          float64
	FontSizeSender    float64
	FontSizeAddress   float64
	LineHeight        float64
	LineHeightAddress float64
	AddressSectionX   float64
	AddressSectionY   float64
	AddressSectionW   float64
	DateY             float64
	Margins           float64
	DatePrefix        string
	Date              string
	Sender            []string
	// Pointers so that we can unset a field that was specified in a more global config
	SenderName *string
	Signature  *string
}

func (c Config) GetSenderNameOrEmpty() string {
	if c.SenderName == nil {
		return ""
	} else {
		return *c.SenderName
	}
}

func (c Config) GetSignatureOrEmpty() string {
	if c.Signature == nil {
		return ""
	} else {
		return *c.Signature
	}
}

var defaultConfig = Config{
	FontName:          "dejavu",
	FontSize:          12,
	FontSizeSender:    7,
	FontSizeAddress:   10,
	LineHeight:        8,
	LineHeightAddress: 6,
	AddressSectionX:   25,
	AddressSectionY:   50,
	AddressSectionW:   70,
	DateY:             100,
	Margins:           25,
	Date:              time.Now().Format("02.01.2006"),
	Sender:            []string{},
}

func printConfiguration(config Config) (string, error) {
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	} else {
		return string(bytes[:]), nil
	}
}

func loadConfigFromFile(configPath string, dest *Config) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			//goland:noinspection GoUnhandledErrorResult
			return errors.New(fmt.Sprintf("Could not read file %s: %s\n", configPath, err))
		}
		return nil
	}
	err = json.Unmarshal(data, &dest)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not parse file %s as json: %s\n", configPath, err))
	} else {
		return nil
	}
}

func loadDefaultConfig(customConfigFilePath string) (Config, error) {
	result := defaultConfig
	var err error
	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "linux" {
		err = loadConfigFromFile("/etc/left/defaults.json", &result)
		if err != nil {
			return result, err
		}
	}
	userDir, err := os.UserConfigDir()
	if err != nil {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintf(os.Stderr, "Could not read user config: %s\n", err)
	} else {
		err = loadConfigFromFile(path.Join(userDir, "left", "defaults.json"), &result)
		if err != nil {
			return result, err
		}
	}
	if customConfigFilePath != "" {
		err = loadConfigFromFile(customConfigFilePath, &result)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func createEmptyLetter(config Config) (string, error) {
	conf, err := printConfiguration(config)
	if err != nil {
		return "", err
	}
	result := "You can put random notes here. Anything before the first section will be ignored.\n"
	result += "Config sections are started with a line that begins with //\n"
	result += "// config\n"
	result += conf + "\n"
	result += "// address\n"
	result += "Name\n"
	result += "Street\n"
	result += "City\n"
	result += "// subject\n"
	result += "Add your subject here. This section must not have more than one line.\n"
	result += "// body\n"
	result += "Dear sir or madam,\n"
	result += "\n"
	result += "\n"
	result += "\n"
	result += "Kind regards,\n"
	return result, nil
}
