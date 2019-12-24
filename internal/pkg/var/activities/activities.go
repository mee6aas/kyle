package activities

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

const rootDirPath = "/act/rsc/asdf"

// AddFromTarGz adds an activity at the specified Gzip with the specified username and activity name to collection.
// The file `actTarGzPath` must be gzipped tarball.
func AddFromTarGz(actName string, actTarGzPath string) (e error) {
	trg := filepath.Join(rootDirPath, actName)

	tgz := archiver.NewTarGz()
	tgz.ImplicitTopLevelFolder = true
	if e = tgz.Unarchive(actTarGzPath, trg); e != nil {
		e = errors.Wrapf(e, "Failed to unarchive %s", actTarGzPath)
		return
	}

	return
}

// AddFromHTTP adds an activity at the specified HTTP address with specified username and activity name to collection.
// The address `actAddr` must be http address and must response gzipped tarball as the attachment.
func AddFromHTTP(actName string, actAddr string) (e error) {
	var (
		res *http.Response
		trg *os.File
	)

	if res, e = http.Get(actAddr); e != nil {
		e = errors.Wrapf(e, "Failed to GET %s", actAddr)
		return
	}
	defer res.Body.Close()

	if trg, e = ioutil.TempFile("", ""); e != nil {
		e = errors.Wrap(e, "Failed to create temporal directory")
		return
	}
	defer os.Remove(trg.Name())

	if _, e = io.Copy(trg, res.Body); e != nil {
		e = errors.Wrapf(e, "Failed to write response from %s", actAddr)
		return
	}

	if e = AddFromTarGz(actName, trg.Name()); e != nil {
		return
	}

	return
}
