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
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	dbg "github.com/visionmedia/go-debug"
)

var (
	debug = dbg.Debug("TEST")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune(" abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type SnapperSuite struct {
	suite.Suite
}

func (suite *SnapperSuite) TestSnapper() {
	rand := randSeq(1024)

	infile, err := ioutil.TempFile("", "")
	assert.NoError(suite.T(), err)
	infile.WriteString(rand)
	assert.NoError(suite.T(), infile.Sync())

	n := infile.Name()
	infile.Close()
	infile, _ = os.Open(n)
	debug("Infile: %s", infile.Name())
	outfile, err := ioutil.TempFile("", "")
	debug("Outfile: %s", outfile.Name())
	assert.NoError(suite.T(), err)
	testCases := []struct {
		desc   string
		input  io.Reader
		output io.ReadWriter
	}{
		{
			desc:   "With Buffers",
			input:  bytes.NewBufferString(rand),
			output: bytes.NewBuffer(nil),
		},
		{
			desc:   "With Files",
			input:  infile,
			output: outfile,
		},
	}
	for _, tC := range testCases {
		suite.T().Run(tC.desc, func(t *testing.T) {
			s := NewSnapper(InFile(tC.input), OutFile(tC.output), Mode(SNAP))
			assert.NoError(t, s.Do())
			result := bytes.NewBuffer(nil)
			if seeker, ok := tC.output.(io.ReadWriteSeeker); ok {
				seeker.Seek(0, 0)
			}
			unsnap := NewSnapper(InFile(tC.output), OutFile(result), Mode(UNSNAP))
			assert.NoError(t, unsnap.Do())
			assert.EqualValues(t, rand, result.String())
		})
	}
}

func (suite *SnapperSuite) TestLackingParams() {
	testCases := []struct {
		desc string
		opts []Option
	}{
		{
			desc: "No Infile",
			opts: []Option{OutFile(bytes.NewBuffer(nil))},
		},
		{
			desc: "No Outfile",
			opts: []Option{InFile(bytes.NewBufferString("test"))},
		},
		{
			desc: "No mode",
			opts: []Option{InFile(bytes.NewBufferString("test")), OutFile(bytes.NewBuffer(nil)), Mode(3)},
		},
	}
	for _, tC := range testCases {
		suite.T().Run(tC.desc, func(t *testing.T) {
			s := NewSnapper(tC.opts...)
			assert.Error(t, s.Do())
		})
	}
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(SnapperSuite))
}
