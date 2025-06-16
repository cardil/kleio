package storage

import (
	"archive/zip"
	"errors"
	"io"
	"os"
)

type FileReader func() io.ReadCloser

type Artifacts map[string]FileReader

func (a Artifacts) ZipReader() (*ZipStream, error) {
	tmp, err := os.CreateTemp("", "artifacts-*.zip")
	if err != nil {
		return nil, err
	}
	remove := true
	defer func(name string) {
		if remove {
			_ = os.Remove(name)
		}
	}(tmp.Name())
	zw := zip.NewWriter(tmp)
	for name, readerFn := range a {
		reader := readerFn()
		var zWriter io.Writer
		if zWriter, err = zw.Create(name); err != nil {
			if cerr := zw.Close(); cerr != nil {
				return nil, errors.Join(err, cerr)
			}
			return nil, err
		}
		if _, err = io.Copy(zWriter, reader); err != nil {
			if cerr := zw.Close(); cerr != nil {
				return nil, errors.Join(err, cerr)
			}
			return nil, err
		}
		if err = reader.Close(); err != nil {
			if cerr := zw.Close(); cerr != nil {
				return nil, errors.Join(err, cerr)
			}
			return nil, err
		}
	}
	if err = zw.Close(); err != nil {
		return nil, err
	}
	var st os.FileInfo
	if st, err = tmp.Stat(); err != nil {
		return nil, err
	}
	if _, err = tmp.Seek(0, 0); err != nil {
		return nil, err
	}
	stream := &ZipStream{
		ReadCloser: &tmpfileReader{tmp},
		Size:       st.Size(),
	}
	remove = false
	return stream, nil
}

type ZipStream struct {
	io.ReadCloser
	Size int64
}

type tmpfileReader struct {
	*os.File
}

func (t *tmpfileReader) Close() (err error) {
	defer func() {
		if cerr := os.Remove(t.Name()); cerr != nil {
			if err != nil {
				err = errors.Join(err, cerr)
			} else {
				err = cerr
			}
		}
	}()
	err = t.File.Close()
	if err != nil {
		return
	}
	return
}
