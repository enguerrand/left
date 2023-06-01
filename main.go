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
	"flag"
	"fmt"
	"log"
	"os"
)

func abort(message string, invocationError bool) {
	printError(message + "\n")
	if invocationError {
		printUsage()
	}
	os.Exit(1)
}

func printError(message string) {
	_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func printUsage() {
	fmt.Println("left - generates letter from txt file")
	fmt.Println("")
	fmt.Println("Usage: left OPTIONS | FILE")
	fmt.Println("")
	fmt.Println("If a FILE argument is provided the file is used as an input.txt to generate a PDF formatted letter.")
	fmt.Println("The text file is expected to consist of two sections, delimited by a line that only contains three")
	fmt.Println("equal signs. (===)")
	fmt.Println("The first section contains the letter configuration, formatted in json (Also see OPTIONS).")
	fmt.Println("The second section contains the letter content.")
	fmt.Println("")
	fmt.Println("Otherwise the following OPTIONS are available:")
	fmt.Println("  -help")
	fmt.Println("        Prints this help")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printUsage
	version := flag.Bool("version", false, "ignore all other arguments, print the left version and exit")
	dumpConfig := flag.Bool("dump-config", false, "dumps the standard config to stdout")
	customConfig := flag.String("config", "", "custom config file to read from after loading configuration defaults")
	create := flag.Bool("create", false, "prints a template for a new letter to stdout")

	flag.Parse()

	if *version {
		fmt.Println(Version())
		os.Exit(0)
	}

	loadedDefaultConfig, err := loadDefaultConfig(*customConfig)
	if err != nil {
		abort(err.Error(), false)
	}
	remainingArgs := os.Args[len(os.Args)-flag.NArg():]
	if *dumpConfig && *create {
		abort("flags -dump-config and -create are mutually exclusive!", true)
	} else if *create && len(remainingArgs) > 0 {
		abort("flag -create is incompatible with positional arguments!", true)
	} else if *create {
		emptyLetter, err := createEmptyLetter(loadedDefaultConfig)
		if err == nil {
			fmt.Println(emptyLetter)
		}
	} else if *dumpConfig {
		configDump, err := printConfiguration(loadedDefaultConfig)
		if err == nil {
			fmt.Println(configDump)
		}
	} else {
		// Consume all the flags that were parsed as flags.
		if len(remainingArgs) == 0 {
			abort("Missing arguments.", true)
		}
		inputFile := remainingArgs[0]
		err = render(inputFile, loadedDefaultConfig)
	}
	if err != nil {
		abort(err.Error(), false)
	}
}
