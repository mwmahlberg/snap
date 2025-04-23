// Copyright ©2016-2022 Markus W Mahlberg <markus@mahlberg.io>
//
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

package pkg

import (
	"bufio"
	"fmt"
	"io"

	"github.com/golang/snappy"
	dbg "github.com/visionmedia/go-debug"
)

type Operation int

const (
	SNAP Operation = iota
	UNSNAP
)

type Snapper struct {
	input  io.Reader
	output io.Writer
	keep   bool
	op     Operation
}

type Option func(*Snapper)

func InFile(in io.Reader) Option {
	return func(s *Snapper) {
		s.input = in
	}
}

func OutFile(out io.Writer) Option {
	return func(s *Snapper) {
		s.output = out
	}
}

func Mode(op Operation) Option {
	return func(s *Snapper) {
		s.op = op
	}
}

func NewSnapper(options ...Option) *Snapper {
	snapper := &Snapper{}

	for _, opt := range options {
		opt(snapper)
	}

	return snapper
}

func (s *Snapper) Keep() bool {
	return s.keep
}

type BufferedWriter interface {
	io.Writer
	Flush() error
}

func (s *Snapper) Do() error {

	if s.input == nil {
		return fmt.Errorf("no input given")
	} else if s.output == nil {
		return fmt.Errorf("no output given")
	}
	var debug = dbg.Debug("OP")

	var in io.Reader
	var out BufferedWriter

	debug("Mode identifier: %d", s.op)
	switch s.op {

	case SNAP:
		debug("Setting up streams for compression")
		in = bufio.NewReader(s.input)
		out = snappy.NewBufferedWriter(s.output)

	case UNSNAP:
		debug("Setting up streams for decompression")
		in = snappy.NewReader(bufio.NewReaderSize(s.input, 32768))
		out = bufio.NewWriter(s.output)
	default:
		return fmt.Errorf("unknown operation identifier: %d", s.op)
	}

	debug("Deferring flush of output")
	defer out.Flush()

	debug("Copying data")
	w, err := io.Copy(out, in)
	debug("Wrote %d bytes", w)
	return err
}
