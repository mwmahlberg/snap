package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mwmahlberg/snap/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
}

func (suite *MainTestSuite) TestIsStdin() {
	assert.False(suite.T(), isStdin())
}

func (suite *MainTestSuite) TestCreateFile() {
	d, err := ioutil.TempDir("", "snap-test")
	assert.NoError(suite.T(), err, "Error creating tempdir: %s", err)
	defer os.RemoveAll(d)

	f, err := createFile(d + "snap-test.txt")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), f)
	defer f.Close()
}

func (suite *MainTestSuite) TestCreateUnsnapFromInfile() {
	d, err := ioutil.TempDir("", "snap-test")
	assert.NoError(suite.T(), err, "Error creating tempdir: %s", err)
	defer os.RemoveAll(d)

	in, err := ioutil.TempFile("", "*.txt.sz")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), in)
	defer os.Remove(in.Name())

	assert.True(suite.T(), strings.HasSuffix(in.Name(), ".txt.sz"))

	o, err := createUnsnapFromInfile(in)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), o)
	defer os.Remove(o.Name())

	assert.False(suite.T(), strings.HasSuffix(o.Name(), ".txt.sz"))
	assert.True(suite.T(), strings.HasSuffix(o.Name(), ".txt"))
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func TestWithInfile(t *testing.T) {

	testCases := []struct {
		desc string
		cfg  config
	}{
		{
			desc: "Compression",
			cfg: config{
				Suffix: ".sz",
			},
		}, {
			desc: "Compression and keep",
			cfg: config{
				Suffix: ".sz",
				Keep:   true,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			f, err := ioutil.TempFile("", "*.txt")
			assert.NoError(t, err)
			assert.NotNil(t, f)
			defer os.Remove(f.Name())
			defer f.Close()

			tC.cfg.InFile = f

			opts, err := setup(tC.cfg)
			fmt.Printf("%+v\n", tC.cfg)
			assert.NoError(t, err)
			assert.NotEmpty(t, opts)
			s := pkg.NewSnapper(opts...)
			assert.NotNil(t, s)
			// in := tC.cfg.InFile.Name()
			// out := s.InFile().Name()
			// assert.EqualValues(t, in, out)
			// assert.EqualValues(t, tC.cfg.Keep, s.Keep())
		})
	}
}
