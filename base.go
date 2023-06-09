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

func MapStrings(input []string, mapper func(string) string) []string {
	output := make([]string, len(input))
	for i := 0; i < len(input); i++ {
		output[i] = mapper(input[i])
	}
	return output
}

func Contains[T comparable](array []T, element T) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}
