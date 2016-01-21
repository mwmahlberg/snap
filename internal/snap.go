// Copyright ©2016 Markus W Mahlberg <markus@mahlberg.io>
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

package internal

import (
	"bufio"
	"io"
	"os"

	"github.com/golang/snappy"
	dbg "github.com/tj/go-debug"
)

type Snapper struct {
	inFileName  string
	outFileName string
	inFile      *os.File
	outFile     *os.File
	inBuf       *bufio.Reader
	outBuf      *bufio.Writer
}

func NewSnapper(infile, outfile string) (*Snapper, error) {
	
	snapper := &Snapper{inFileName: infile, outFileName: outfile}

	if in, err := os.OpenFile(infile, os.O_RDONLY, 0600); err == nil {
		snapper.inFile = in
	} else {
		return nil, err
	}

	snapper.inBuf = bufio.NewReader(snapper.inFile)

	inFi, _ := snapper.inFile.Stat()
	if out, err := os.OpenFile(outfile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, inFi.Mode()); err == nil {
		snapper.outFile = out
	} else {
		return nil, err
	}
	snapper.outBuf = bufio.NewWriter(snapper.outFile)

	return snapper, nil
}

func (s *Snapper) Snap() error {
	debug := dbg.Debug("SNAP")
	defer s.inFile.Close()
	
	snap := bufio.NewWriter(snappy.NewWriter(s.outFile))
	defer s.outFile.Close()
	defer snap.Flush()

	if w, err := io.Copy(snap, s.inBuf); err != nil {
		debug("Error compressing file after %d bytes: %v", w, err)
		return err
	}

	return nil
}

func (s *Snapper) Unsnap() error {
	debug := dbg.Debug("UNSNAP")

	defer s.inFile.Close()
	defer s.outFile.Close()
	defer s.outBuf.Flush()

	usnap := snappy.NewReader(bufio.NewReader(s.inFile))

	if w, err := io.Copy(s.outBuf, usnap); err != nil {
		debug("Error decompressing file after %d bytes: %v", w, err)
		return err
	}
	return nil
}
