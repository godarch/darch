// Copyright 2015 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package indent indents lines of text.
package indent

import (
	"bytes"
	"io"
	"strings"
)

// String returns s with each line in s prefixed by indent.
func String(indent, s string) string {
	if indent == "" || s == "" {
		return s
	}
	lines := strings.SplitAfter(s, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(append([]string{""}, lines...), indent)
}

// Bytes returns b with each line in b prefixed by indent.
func Bytes(indent, b []byte) []byte {
	if len(indent) == 0 || len(b) == 0 {
		return b
	}
	lines := bytes.SplitAfter(b, []byte{'\n'})
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return bytes.Join(append([][]byte{[]byte{}}, lines...), indent)
}

// NewWriter returns an io.Writer that prefixes the lines written to it with
// indent and then writes them to w.  The writer returns the number of bytes
// written to the underlying Writer.
func NewWriter(w io.Writer, indent string) io.Writer {
	if indent == "" {
		return w
	}
	return &iw{
		w:      w,
		prefix: []byte(indent),
	}
}

type iw struct {
	w       io.Writer
	prefix  []byte
	partial bool // true if next line's indent already written
}

// Write implements io.Writer.
func (w *iw) Write(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	lines := bytes.SplitAfter(buf, []byte{'\n'})
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	if !w.partial {
		lines = append([][]byte{[]byte{}}, lines...)
	}
	joined := bytes.Join(lines, w.prefix)
	w.partial = joined[len(joined)-1] != '\n'

	n, err := w.w.Write(joined)
	if err != nil {
		return actualWrittenSize(n, len(w.prefix), lines), err
	}

	return len(buf), nil
}

func actualWrittenSize(underlay, prefix int, lines [][]byte) int {
	actual := 0
	remain := underlay
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		addition := remain - prefix
		if addition <= 0 {
			return actual
		}

		if addition <= len(line) {
			return actual + addition
		}

		actual += len(line)
		remain -= prefix + len(line)
	}

	return actual
}
