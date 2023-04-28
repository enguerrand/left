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
	"bufio"
	"encoding/json"
	"errors"
	"github.com/go-pdf/fpdf"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var sectionSeparationRegex = regexp.MustCompile("^//.*")

type LetterSection int

const (
	Initial       LetterSection = iota
	Configuration LetterSection = iota
	Address       LetterSection = iota
	Subject       LetterSection = iota
	Body          LetterSection = iota
)

func render(inputFile string, defaultConfig Config) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	var text []string
	var configJson string
	var subject = ""
	var recipient []string
	var bodyReached = false
	var multiLineSubject = false
	file, jsonReadError := os.Open(inputFile)
	if jsonReadError != nil {
		return jsonReadError
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()
	scanner := bufio.NewScanner(file)
	sectionIndex := 0

	for scanner.Scan() {
		line := scanner.Text()
		if sectionSeparationRegex.MatchString(line) {
			sectionIndex++
			if LetterSection(sectionIndex) == Body {
				bodyReached = true
			}
			continue
		}
		switch LetterSection(sectionIndex) {
		case Initial:
			// Tolerate freestyle text before the config section
			continue
		case Configuration:
			configJson = configJson + line
		case Subject:
			if subject != "" {
				/*
				 Don't error out just yet in order to produce meaningful error messages. If we end up finding
				 the correct count of config sections we will complain about the multiline subject.
				 If however we detect missing config sections the multiline subject is just a symptom and we should
				 really complain about missing sections.
				*/
				multiLineSubject = true
			}
			subject = tr(line)
		case Address:
			recipient = append(recipient, tr(line))
		case Body:
			fallthrough // tolerate config separator in body
		default:
			text = append(text, tr(scanner.Text()))
		}
	}
	for scanner.Scan() {
		text = append(text, tr(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if !bodyReached {
		return errors.New("letters MUST have exactly four sections: config, address, subject, body (in this order), initiated by lines starting with //")
	} else if multiLineSubject {
		return errors.New("the subject section (third section) must only contain one single line")
	}
	utf8Config := defaultConfig
	jsonParseError := json.Unmarshal([]byte(configJson), &utf8Config)
	if jsonParseError != nil {
		return jsonParseError
	}
	trSenderName := tr(utf8Config.GetSenderNameOrEmpty())
	trSignature := tr(utf8Config.GetSignatureOrEmpty())
	config := Config{
		utf8Config.FontName,
		utf8Config.FontSize,
		utf8Config.FontSizeSender,
		utf8Config.FontSizeAddress,
		utf8Config.LineHeight,
		utf8Config.LineHeightAddress,
		utf8Config.AddressSectionX,
		utf8Config.AddressSectionY,
		utf8Config.AddressSectionW,
		utf8Config.DateY,
		utf8Config.Margins,
		tr(utf8Config.DatePrefix),
		tr(utf8Config.Date),
		mapStrings(utf8Config.Sender, tr),
		&trSenderName,
		&trSignature,
	}

	pdf.AddPage()
	pdf.SetMargins(config.Margins, 20, config.Margins)

	// Sender
	pdf.SetXY(config.AddressSectionX, config.AddressSectionY)
	pdf.SetFont(config.FontName, "", config.FontSizeSender)
	pdf.MultiCell(config.AddressSectionW, config.LineHeightAddress, strings.Join(config.Sender, ", "), "B", "L", false)

	// Address
	pdf.SetFont(config.FontName, "", config.FontSizeAddress)
	for i := 0; i < len(recipient); i++ {
		pdf.SetX(config.AddressSectionX)
		pdf.MultiCell(config.AddressSectionW, config.LineHeightAddress, recipient[i], "", "L", false)
	}

	pdf.SetFont(config.FontName, "", config.FontSize)

	// Date
	pdf.SetXY(config.Margins, config.DateY)
	pdf.MultiCell(0, config.LineHeight, config.DatePrefix+config.Date, "", "R", false)

	// Subject
	pdf.SetFont(config.FontName, "B", config.FontSize)
	pdf.SetXY(config.Margins, config.DateY+config.LineHeight)
	pdf.MultiCell(0, config.LineHeight, subject, "", "L", false)
	pdf.SetFont(config.FontName, "", config.FontSize)

	pdf.Ln(config.LineHeight)

	// Text
	for i := 0; i < len(text); i++ {
		pdf.SetX(config.Margins)
		pdf.MultiCell(0, config.LineHeight, text[i], "", "L", false)
	}

	if config.GetSignatureOrEmpty() != "" {
		var opt fpdf.ImageOptions
		opt.ImageType = "jpg"
		pdf.ImageOptions(config.GetSignatureOrEmpty(), pdf.GetX(), pdf.GetY(), 0, 0, false, opt, 0, "")
	}
	pdf.Ln(config.LineHeight * 3)
	pdf.MultiCell(0, config.LineHeight, config.GetSenderNameOrEmpty(), "", "L", false)

	pdfErr := pdf.OutputFileAndClose(strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".pdf")
	return pdfErr
}
