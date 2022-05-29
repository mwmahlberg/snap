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

package pkg

import (
	"bufio"
	"io"

	"github.com/golang/snappy"
	dbg "github.com/visionmedia/go-debug"
)

type Snapper struct {
	in  io.Reader
	out io.Writer
}

func NewSnapper(in io.Reader, out io.Writer) *Snapper {

	snapper := &Snapper{out: out, in: in}

	return snapper
}

func (s *Snapper) Snap() error {
	debug := dbg.Debug("SNAP")

	snap := snappy.NewBufferedWriter(s.out)
	defer snap.Flush()
	if w, err := io.Copy(snap, s.in); err != nil {
		debug("Error compressing file after %d bytes: %v", w, err)
		return err
	} else {
		debug("Wrote %d bytes", w)
	}

	return nil
}

func (s *Snapper) Unsnap() error {
	debug := dbg.Debug("UNSNAP")

	usnap := snappy.NewReader(bufio.NewReader(s.in))

	if w, err := io.Copy(s.out, usnap); err != nil {
		debug("Error decompressing file after %d bytes: %v", w, err)
		return err
	}
	return nil
}
