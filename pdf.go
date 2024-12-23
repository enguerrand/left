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
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-pdf/fpdf"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var sectionSeparationRegex = regexp.MustCompile("^//.*")

//go:embed fonts/*
var fontsDir embed.FS

type LetterSection int

const (
	Initial       LetterSection = iota
	Configuration LetterSection = iota
	Address       LetterSection = iota
	Subject       LetterSection = iota
	Body          LetterSection = iota
)

func addEmbeddedFont(pdf *fpdf.Fpdf, family string) {
	regular, _ := fontsDir.ReadFile(fmt.Sprintf("fonts/%s/regular.ttf", family))
	for _, style := range []string{"", "I"} {
		pdf.AddUTF8FontFromBytes(family, style, regular)
	}
	bold, _ := fontsDir.ReadFile(fmt.Sprintf("fonts/%s/bold.ttf", family))
	for _, style := range []string{"B", "BI"} {
		pdf.AddUTF8FontFromBytes(family, style, bold)
	}
}

func addExternalFont(pdf *fpdf.Fpdf, fontImport FontImport) {
	pdf.SetFontLocation(fontImport.Directory)
	for _, styleString := range []string{"", "B", "I", "BI"} {
		if fontImport.FontFileNameBold != "" && (styleString == "B" || styleString == "BI") {
			pdf.AddUTF8Font(fontImport.Name, styleString, fontImport.FontFileNameBold)
		} else {
			pdf.AddUTF8Font(fontImport.Name, styleString, fontImport.FontFileName)
		}
	}
}

func render(inputFile string, defaultConfig Config) error {
	pdf := fpdf.New("P", "mm", "A4", "")

	utf8Fonts := []string{"dejavusanscondensed", "freeserif"}
	for _, family := range utf8Fonts {
		addEmbeddedFont(pdf, family)
	}

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
			subject = line
		case Address:
			recipient = append(recipient, line)
		case Body:
			fallthrough // tolerate config separator in body
		default:
			text = append(text, scanner.Text())
		}
	}
	for scanner.Scan() {
		text = append(text, scanner.Text())
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

	if utf8Config.FontImport != nil {
		addExternalFont(pdf, *utf8Config.FontImport)
		utf8Fonts = append(utf8Fonts, utf8Config.FontImport.Name)
	}

	fontName := utf8Config.FontName
	var tr func(s string) string
	if Contains(utf8Fonts, strings.ToLower(fontName)) {
		tr = func(s string) string {
			return s
		}
	} else {
		tr = pdf.UnicodeTranslatorFromDescriptor("")
	}

	trSubject := tr(subject)
	trRecipient := MapStrings(recipient, tr)
	trText := MapStrings(text, tr)
	trSenderName := tr(utf8Config.GetSenderNameOrEmpty())
	trSignature := tr(utf8Config.GetSignatureOrEmpty())
	config := Config{
		utf8Config.FontName,
		utf8Config.FontImport,
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
		MapStrings(utf8Config.Sender, tr),
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
	for i := 0; i < len(trRecipient); i++ {
		pdf.SetX(config.AddressSectionX)
		pdf.MultiCell(config.AddressSectionW, config.LineHeightAddress, trRecipient[i], "", "L", false)
	}

	pdf.SetFont(config.FontName, "", config.FontSize)

	// Date
	pdf.SetXY(config.Margins, config.DateY)
	pdf.MultiCell(0, config.LineHeight, config.DatePrefix+config.Date, "", "R", false)

	// Subject
	pdf.SetFont(config.FontName, "B", config.FontSize)
	pdf.SetXY(config.Margins, config.DateY+config.LineHeight)
	pdf.MultiCell(0, config.LineHeight, trSubject, "", "L", false)
	pdf.SetFont(config.FontName, "", config.FontSize)

	pdf.Ln(config.LineHeight)

	// Text
	for i := 0; i < len(trText); i++ {
		pdf.SetX(config.Margins)
		pdf.MultiCell(0, config.LineHeight, trText[i], "", "L", false)
	}

	if config.GetSignatureOrEmpty() != "" {
		var opt fpdf.ImageOptions
		opt.ImageType = "jpg"
		pdf.ImageOptions(config.GetSignatureOrEmpty(), pdf.GetX(), pdf.GetY(), 0, 0, true, opt, 0, "")
	}
	pdf.Ln(config.LineHeight)
	pdf.MultiCell(0, config.LineHeight, config.GetSenderNameOrEmpty(), "", "L", false)

	pdfErr := pdf.OutputFileAndClose(strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".pdf")
	return pdfErr
}
