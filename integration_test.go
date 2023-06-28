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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

/*
This test depends on diff-pdf being present (https://vslavik.github.io/diff-pdf/)
*/
func TestPdfCreation(t *testing.T) {
	resDir := "./test/it/pdf"
	files, err := os.ReadDir(resDir)
	if err != nil {
		t.Errorf("Could not read test resources: %s", err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			outfile := filepath.Join(resDir, file.Name(), "input.pdf")
			reference := filepath.Join(resDir, file.Name(), "reference.pdf")
			_ = os.Remove(outfile)
			inputFile := filepath.Join(resDir, file.Name(), "input.left")
			f := false
			configPaths := []string{}
			index := 0
			for {
				index++
				config := filepath.Join(resDir, file.Name(), fmt.Sprintf("config_%d.json", index))
				if _, err := os.Stat(config); err == nil {
					configPaths = append(configPaths, config)
				} else {
					break
				}
			}
			Run(configPaths, &f, &f, []string{inputFile})

			cmd := exec.Command("diff-pdf", outfile, reference)
			if err := cmd.Run(); err != nil {
				t.Errorf("Created pdf %s did not match the reference: %s. To compare, run:\n\ndiff-pdf --view %s %s\n\n", outfile, reference, outfile, reference)
			} else {
				_ = os.Remove(outfile)
			}
		}
	}
}

func TestCreate(t *testing.T) {
	resDir := "./test/it/create"
	files, err := os.ReadDir(resDir)
	if err != nil {
		t.Errorf("Could not read test resources: %s", err.Error())
	}
	for _, file := range files {
		if file.IsDir() {

			reference := filepath.Join(resDir, file.Name(), "expect.left")
			expected, err := os.ReadFile(reference)
			if err != nil {
				t.Errorf("Failed to load expect.left: " + err.Error())
			}
			configPaths := []string{}
			index := 0
			for {
				index++
				config := filepath.Join(resDir, file.Name(), fmt.Sprintf("config_%d.json", index))
				if _, err := os.Stat(config); err == nil {
					configPaths = append(configPaths, config)
				} else {
					break
				}
			}
			loadedDefaultConfig, err := loadDefaultConfig(configPaths)
			if err != nil {
				t.Errorf("Failed to load config: " + err.Error())
			}
			emptyLetter, err := createEmptyLetter(loadedDefaultConfig)
			if err != nil {
				t.Errorf("Failed to create empty letter: " + err.Error())
			}
			if string(expected) != emptyLetter {
				t.Errorf(fmt.Sprintf("Created empty letter does not the match the result. Expected:\n%s\n\n Created:\n%s", expected, emptyLetter))
			}
		}
	}
}

func TestDump(t *testing.T) {
	resDir := "./test/it/dump"
	files, err := os.ReadDir(resDir)
	if err != nil {
		t.Errorf("Could not read test resources: %s", err.Error())
	}
	for _, file := range files {
		if file.IsDir() {

			reference := filepath.Join(resDir, file.Name(), "expect.json")
			expected, err := os.ReadFile(reference)
			if err != nil {
				t.Errorf("Failed to load expect.json: " + err.Error())
			}
			configPaths := []string{}
			index := 0
			for {
				index++
				config := filepath.Join(resDir, file.Name(), fmt.Sprintf("config_%d.json", index))
				if _, err := os.Stat(config); err == nil {
					configPaths = append(configPaths, config)
				} else {
					break
				}
			}
			loadedDefaultConfig, err := loadDefaultConfig(configPaths)
			if err != nil {
				t.Errorf("Failed to load config: " + err.Error())
			}
			emptyLetter, err := printConfiguration(loadedDefaultConfig)
			if err != nil {
				t.Errorf("Failed to dump configuration: " + err.Error())
			}
			if string(expected) != emptyLetter {
				t.Errorf(fmt.Sprintf("Dumped configuration does not the match the result. Expected:\n%s \n\n Dumped:\n%s", expected, emptyLetter))
			}
		}
	}
}
