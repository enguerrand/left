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
	"reflect"
	"testing"
)

func TestReadMinimalConfigFromFile(t *testing.T) {
	read := Config{}
	var fontImport *FontImport = nil
	var nilSlice []string
	var nilStringPtr *string
	_ = loadConfigFromFile("./test/config/valid_config_minimal.json", &read)
	AssertEquals(t, read.FontName, "", "FontName")
	AssertEquals(t, read.FontImport, fontImport, "FontImport")
	AssertEquals(t, read.FontSize, float64(0), "FontSize")
	AssertEquals(t, read.FontSizeSender, float64(0), "FontSizeSender")
	AssertEquals(t, read.FontSizeAddress, float64(0), "FontSizeAddress")
	AssertEquals(t, read.LineHeight, float64(0), "LineHeight")
	AssertEquals(t, read.LineHeightAddress, float64(0), "LineHeightAddress")
	AssertEquals(t, read.AddressSectionX, float64(0), "AddressSectionX")
	AssertEquals(t, read.AddressSectionY, float64(0), "AddressSectionY")
	AssertEquals(t, read.AddressSectionW, float64(0), "AddressSectionW")
	AssertEquals(t, read.DateY, float64(0), "DateY")
	AssertEquals(t, read.Margins, float64(0), "Margins")
	AssertEquals(t, read.DatePrefix, "", "DatePrefix")
	AssertEquals(t, read.Date, "", "Date")
	AssertStringSliceEquals(t, nilSlice, read.Sender, "Sender")
	AssertEquals(t, nilStringPtr, read.SenderName, "SenderName")
	AssertEquals(t, nilStringPtr, read.Signature, "Signature")
}

func TestReadFullConfigFromFile(t *testing.T) {
	read := Config{}
	_ = loadConfigFromFile("./test/config/valid_config_full.json", &read)
	AssertEquals(t, read.FontName, "someFontName", "FontName")
	AssertEquals(t, read.FontImport.Name, "myfont", "FontImport.Name")
	AssertEquals(t, read.FontImport.Directory, "/usr/share/fonts/myfont", "FontImport.Directory")
	AssertEquals(t, read.FontImport.FontFileName, "MyFont-Condensed.ttf", "FontImport.FontFileName")
	AssertEquals(t, read.FontSize, float64(42), "FontSize")
	AssertEquals(t, read.FontSizeSender, float64(43), "FontSizeSender")
	AssertEquals(t, read.FontSizeAddress, float64(44), "FontSizeAddress")
	AssertEquals(t, read.LineHeight, float64(45), "LineHeight")
	AssertEquals(t, read.LineHeightAddress, float64(46), "LineHeightAddress")
	AssertEquals(t, read.AddressSectionX, float64(47), "AddressSectionX")
	AssertEquals(t, read.AddressSectionY, float64(48), "AddressSectionY")
	AssertEquals(t, read.AddressSectionW, float64(49), "AddressSectionW")
	AssertEquals(t, read.DateY, float64(50), "DateY")
	AssertEquals(t, read.Margins, float64(51), "Margins")
	AssertEquals(t, read.DatePrefix, "My Hometown, ", "DatePrefix")
	AssertEquals(t, read.Date, "24/05/2023", "Date")
	AssertEquals(t, read.Sender[0], "Darth Vader", "read.Sender[0]")
	AssertEquals(t, read.Sender[1], "Palace District with Special Chars äüößéç", "read.Sender[1]")
	AssertEquals(t, read.Sender[2], "Coruscant", "read.Sender[2]")
	AssertEquals(t, *read.SenderName, "Darth Vader", "read.SenderName")
	AssertEquals(t, *read.Signature, "/home/dvader/documents/Signature.jpg", "read.Signature")
}

func AssertEquals(t *testing.T, got any, want any, description string) {
	if got != want {
		t.Errorf("%s: got %q, wanted %q", description, got, want)
	}
}
func AssertStringSliceEquals(t *testing.T, got []string, want []string, description string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s: got %q, wanted %q", description, got, want)
	}
}
