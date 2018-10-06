// Code generated by vfsgen; DO NOT EDIT.

package assets

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// assets statically implements the virtual filesystem provided to vfsgen.
var assets = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(2018, 10, 6, 11, 12, 0, 346162909, time.UTC),
		},
		"/circuit": &vfsgen۰DirInfo{
			name:    "circuit",
			modTime: time.Date(2018, 10, 6, 11, 11, 43, 291283315, time.UTC),
		},
		"/circuit/NOTE.txt": &vfsgen۰CompressedFileInfo{
			name:             "NOTE.txt",
			modTime:          time.Date(2018, 10, 6, 11, 11, 43, 288051929, time.UTC),
			uncompressedSize: 237,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\xce\xb1\x4e\x03\x31\x10\x04\xd0\xde\x5f\x31\xfc\x80\xaf\xbf\x8e\x86\x50\x22\x1a\x44\x15\x6d\x7c\x73\xc4\x62\xe3\xb5\xbc\x36\x28\x20\xfe\x1d\x1d\x1c\x74\xbb\x23\x8d\xe6\xdd\x53\xd5\xf0\x64\x4d\x17\x54\x49\xaf\xf2\x42\xf8\x48\x89\xee\xeb\x50\xbd\x62\x61\x55\xbb\x72\xb9\x09\xe1\x41\x29\x4e\x34\xae\x6c\xe8\x67\xfe\xdc\x6e\xa3\x25\xfa\x8c\x10\x1e\xf7\x07\x87\x66\xa3\x02\x33\x62\x8c\xe1\xf6\x63\x34\xe2\x6e\x94\xd4\xb3\x15\xdf\xd3\x70\x30\x74\xc3\xb9\xf7\xea\xf3\x34\x7d\xb2\xbc\xe5\x66\xe5\xc2\xd2\x8f\x27\x71\x1e\x8b\x5c\xf8\x25\xb5\x46\xd9\xfa\xef\x3c\x79\xee\xf4\x58\xd8\x27\xa9\x79\xda\xdc\xbf\x6a\x3c\xdb\x40\x92\x02\x27\x37\x15\x44\x15\xeb\xff\xdc\x9f\x3f\x22\x84\xef\x00\x00\x00\xff\xff\xac\xa6\x50\x8b\xed\x00\x00\x00"),
		},
		"/circuit/manifest.yaml": &vfsgen۰CompressedFileInfo{
			name:             "manifest.yaml",
			modTime:          time.Date(2018, 10, 6, 11, 40, 28, 598048499, time.UTC),
			uncompressedSize: 474,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x51\xbb\x6e\xc2\x40\x10\xec\xef\x2b\x56\xa2\xe6\x4c\xda\xeb\xa2\x48\x79\x34\x08\x01\x51\xea\xc3\x5e\xfb\x36\xdc\x4b\xb7\x6b\x90\x41\xfc\x7b\xe4\x43\xd0\xa7\x9e\xc7\xce\xcc\x2e\x60\x27\x85\x8e\xc8\x90\x6d\x7b\xb4\x03\x42\xb0\x91\x7a\x64\x81\x9e\x3c\x2a\x15\x6d\x40\x03\xd7\xab\xde\xdc\xf1\xb5\x0d\x78\xbb\xa9\x0e\xb9\x2d\x94\x85\x52\x34\xf0\x89\xde\xa7\x9f\x54\x7c\xf7\x74\xe9\x53\x81\x0e\xb3\x4f\x13\xc5\x01\x98\x42\xf6\x08\x6f\x0b\xe8\xc7\xd8\xce\x22\x86\x33\x89\x83\x36\x45\x1e\x43\xb5\x81\xec\x6d\x7c\xe2\x60\x73\xd6\xca\x8e\xe2\x52\x31\xb0\xe7\x71\x4a\xec\x08\xbe\xd9\x51\x52\xb9\xa4\x5f\x6c\x65\x63\x07\x34\xe0\x44\x32\x9b\xa6\x19\x48\xdc\x78\xd0\x6d\x0a\xcd\x83\x5e\xd9\x8d\x9b\xd3\x2d\xcf\x73\xbc\x87\x72\x8b\x39\xfd\x53\xa9\x16\xb0\x45\x8f\x96\x51\xa9\x13\x16\xae\xc5\x5f\xf4\x4a\xaf\x66\xd3\x13\x75\x58\xf6\x53\x46\x03\x7b\x2c\xc5\xf6\xa9\x04\x55\xee\xfc\x75\x12\x34\xf0\x15\x49\xc8\x7a\xa8\x4b\x04\x8c\x62\x6b\xcd\x79\x27\x71\x08\x3d\x15\x16\xf8\xd0\xaf\x1a\xd4\x89\x98\x0e\xe4\x49\x26\x03\x79\x3c\x78\x6a\x15\x8b\x2d\xb2\xab\x8b\x1b\x90\xc7\x05\x2d\xbd\xba\x50\x7e\x27\x5f\xdf\xc2\x46\x2d\xa1\x66\xd6\x17\xca\x7f\x01\x00\x00\xff\xff\xda\xc0\xb9\x3e\xda\x01\x00\x00"),
		},
		"/circuit/terraform.tf": &vfsgen۰CompressedFileInfo{
			name:             "terraform.tf",
			modTime:          time.Date(2018, 10, 6, 11, 11, 43, 290918022, time.UTC),
			uncompressedSize: 3188,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x56\xdb\x6e\xdb\x38\x13\xbe\xd7\x53\x10\x4c\x2f\xfe\x1f\x48\x64\x67\x0f\x17\x5d\x20\x17\x69\xea\x64\x8d\x6c\x9c\xc0\x72\x10\x14\x8b\x42\xa0\xa5\xb1\xcc\x5a\x26\x09\x92\x72\xe2\x04\xda\x67\x5f\xf0\x60\x9d\x2c\xbb\x29\xba\xbe\xb2\x38\x33\x1f\x67\xbe\x39\xf1\x04\x45\x20\x37\x34\x01\xf4\x20\x29\x4b\xa8\xc8\x21\xd8\x10\x49\xc9\x3c\x07\x84\x55\x31\x57\x89\xa4\x42\x53\xce\x62\x9a\x62\xf4\x56\x36\xa4\x49\x4e\x81\xe9\x83\xe7\x0a\x12\x09\xba\x23\xd3\xc0\x48\x6d\x13\x9c\xa0\x29\x28\x5e\xc8\x04\xd0\x8d\xe4\x85\x68\x68\x4a\x2f\x88\x33\x23\xc0\xe8\x2d\x40\x28\x85\x05\x29\x72\x8d\x2e\x10\x5e\x42\x9e\xf3\xb3\x67\x2e\xf3\xf4\x4c\x66\x38\x68\x5e\x92\xf3\x84\x18\x97\xf7\x8c\xbe\x11\x41\x18\x10\xa5\xdb\xfa\x9a\x64\x31\x23\x6b\x38\x76\x89\xb1\x68\x98\x00\xdb\x50\xc9\xd9\xda\x04\x3a\x27\x0a\xde\x67\x6f\xc2\x15\x5c\x51\xcd\xe5\x16\x45\xa0\x35\x65\x99\x0a\x5a\x41\xef\xc4\x0e\xb5\x90\xd4\x81\xb6\x61\xb5\x16\xea\x8f\xc1\x40\x69\x49\x57\xa0\x8c\x11\x7c\x04\xf8\xf6\xfb\xcb\xaf\xe1\x3c\xe7\xf3\x30\xe1\x12\xc2\x67\xca\x52\xfe\xac\x42\x06\x7a\x50\xe3\x0e\xbc\x23\x0f\x24\x21\xab\x0c\x90\xea\xf1\x42\x90\x64\x45\xb2\x66\x4c\x27\x08\x5e\xc8\xda\x14\x87\xf9\xff\x2e\x86\x76\x20\x1b\x90\xaa\xca\xc5\x21\x9c\xf3\x70\x18\x0e\xfb\xed\x5f\xa9\xb0\x8e\xc4\xc3\x1d\x13\x4d\x90\x3e\x77\xc2\x57\x2a\x7c\x94\x91\x80\x84\x2e\xb6\x48\x2f\x01\xe5\x84\x65\x05\xc9\x00\xf1\x85\xfd\xf6\x17\x34\xeb\xc6\x6b\xf4\x50\x9e\x72\xcd\xc0\x96\x4d\x20\x24\xdf\xd0\x14\x24\xc2\xe4\xb5\x90\x20\xd7\x5e\xbd\x6a\x07\xa3\xff\xe1\x6d\x43\x64\x58\x1d\x95\xb8\xa1\xe1\x1a\x63\x4f\xcb\x1d\x3b\xcd\x4e\xe3\xd5\xba\x1d\x81\xd3\xae\x9a\xaa\xd6\xab\x8e\x4a\xeb\xf3\xae\x99\x2a\x9f\xe3\x6e\x7b\x61\x0d\x4a\xbb\x48\x0c\xdb\x26\xfc\x1a\xad\xad\x6c\x2f\xdd\xf5\x58\xad\xb4\x3b\xb1\x62\x4d\x32\xe5\x59\xb4\x70\xb5\x63\xbe\xd7\xac\x56\xd9\xf6\x4d\x12\x96\xf2\x75\x6c\xca\x9a\x65\xd8\x8c\x9f\xc5\x82\xbe\x38\xa7\x72\x60\x99\x5e\xa2\x0b\xf4\x5b\x80\x90\x32\x69\x25\x39\xba\x40\x0b\x92\x2b\x53\x07\x85\x10\x20\xab\xef\xde\x88\x95\xe6\xd2\x14\x14\x49\x12\x5e\x30\xdd\x1f\x72\xf7\x67\xfd\x96\x20\x72\x92\xc0\xff\x8c\xff\xbd\x8d\x7f\x8a\xcf\xf0\x29\xc6\xff\x2f\x15\xf9\xf0\xd6\x8a\x22\x74\x31\x18\x06\x8b\xdc\x25\xb7\x4d\x66\x5c\x5f\x6d\xef\xea\xcf\x4f\x68\x5c\x0d\x2b\xde\x2a\xf2\x7b\xdd\x3d\x06\xd1\x4a\x92\x67\x22\xd6\x14\x64\x17\x26\xd2\x84\xa5\x44\xa6\x4d\x3d\xc3\x03\x75\xf6\xb1\xde\x0a\x9b\xd5\xbf\xa6\x11\xb6\x43\x4f\x00\x4b\x55\x6c\x0b\xe2\x6f\xdc\xc7\x01\xfe\xda\x9f\x17\x22\x44\xac\xdc\x12\x8a\x45\x4e\xd8\x77\x13\x53\xd5\x52\x6f\x2e\x4a\x8b\x71\x80\xa4\x1f\xe3\xa7\x2f\x53\xef\x4f\xd2\x8a\xb2\xb4\xc7\xf5\xeb\x82\x25\xe6\x86\x4b\x21\x70\x60\x6a\x79\x55\xf8\x46\xb1\x69\xb8\x40\xf8\xf3\x96\x91\x35\x4d\x6c\x6b\x23\x45\x5f\xed\xa5\x5f\xce\x7b\x1a\xa6\x41\x61\x95\x18\xca\x14\xcd\x96\x5a\xfd\x2c\x8d\x84\x1e\x21\xd1\xec\xd0\x42\xfd\x3c\x47\x4d\xc7\x6d\x45\x39\xf8\x27\x98\xfb\x01\x7e\xc5\xd5\x9a\xab\xcf\x9f\x7a\xc6\x04\x65\x1a\x32\x90\x18\xe1\x7a\x49\xae\xa9\xa9\xbf\xf3\xe1\x70\x38\x74\xdf\xe4\x05\x5d\xa0\x8f\xe6\xd7\xcf\xdc\xc2\x67\xc3\x50\xf8\xbe\x89\xf0\x5d\xde\x4c\x5e\x8f\xb4\xe8\x7f\x50\x83\x3f\x3a\x2e\xba\x1d\x66\x36\x45\x0f\x46\x57\xcd\xa1\xf8\x25\xb3\x9b\x9e\x09\x67\x0c\x1c\x65\xae\xbb\xdb\x18\x9d\x21\xeb\x20\x84\xa4\x6b\x22\xb7\xfb\xb6\x16\xd9\x3f\x10\x0c\xce\x3f\xbf\xd8\x96\x70\x8e\xb8\x77\x89\x4f\x2c\xbe\x7c\x78\x18\x4f\xa2\xf1\xcd\x9f\xb3\x28\x1e\x4f\xa2\xd9\xf4\xf1\x6e\x34\x99\x5d\xce\xc6\xf7\x93\xdb\xd1\x17\xbc\x17\xc9\x5e\x3f\xf8\x68\x98\xd2\xb2\x30\x29\x73\xd2\x15\x6c\xdd\x0e\x45\xf8\xfa\x71\x72\x65\xe0\xa2\xf8\xe9\x7e\x7a\x3b\x9a\xc6\xd3\xc7\xc9\x6c\x7c\x37\xc2\x8d\x15\xe7\x1f\x08\xde\xe4\x04\xcd\x96\x54\xa1\x8c\x1b\x1a\x34\x47\xc9\x92\xb0\x0c\x10\x65\x9a\xa3\xa7\xd1\xa7\x68\x3c\x1b\x19\x90\xf8\x7a\x7a\x7f\x17\x3f\x5c\x5e\xdd\x5e\xde\x8c\xdc\x5d\x87\xa4\xb8\xb9\x73\xf7\x9e\x82\xa5\x93\x34\x9f\x67\xe5\xa0\x7d\xe6\xc9\x2c\x07\xfe\xa0\x23\xae\x1f\x53\x7e\xff\xee\x8f\xee\xa3\xd5\x80\x4f\x0f\xbd\x20\xf6\xc4\x7d\x95\x80\xbf\xa2\xa0\x0c\xfe\x0d\x00\x00\xff\xff\x1b\x43\x80\x5e\x74\x0c\x00\x00"),
		},
		"/circuit/values.hcl": &vfsgen۰CompressedFileInfo{
			name:             "values.hcl",
			modTime:          time.Date(2018, 10, 6, 12, 13, 57, 322316359, time.UTC),
			uncompressedSize: 349,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8f\x51\xaa\xc2\x30\x10\x45\xff\xb3\x8a\x21\xbf\xef\xf5\x75\x05\x6f\x0b\xe2\x0e\xca\x6d\x3b\xc6\x68\x3b\x29\xc9\xa4\x20\x25\x7b\x17\x41\xb1\x45\x41\xff\xe7\xdc\x73\x66\x46\xf4\x68\x07\x26\xcb\x32\xfb\x18\x64\x64\xd1\xa6\x45\xe2\x46\x30\xb2\xa5\xc5\x10\xf5\x7c\x40\x1e\x94\xfe\xc9\x2e\xcb\xdf\x1e\xdd\x19\x8e\x77\x18\xb9\x14\x6b\x8a\x79\x6e\x44\x4e\x21\xc7\x8e\x1b\x17\x43\x9e\x3e\xc2\x55\x74\x96\xea\x9a\xda\xcb\xe3\xea\x97\xb6\x13\xe4\x13\xe9\x91\xe9\x6d\x1c\xfd\x50\x15\x9d\x29\x66\x55\x30\x84\x0e\xea\x83\xbc\xb8\x4f\x98\x20\x8c\xa4\x76\x0b\x28\xdc\xf7\x9f\xae\x45\x10\x97\xe1\xee\xdc\x86\xec\x83\x0a\xdf\x3c\xd7\x00\x00\x00\xff\xff\xbd\x23\xeb\xb9\x5d\x01\x00\x00"),
		},
		"/package": &vfsgen۰FileInfo{
			name:    "package",
			modTime: time.Date(2018, 10, 6, 11, 12, 0, 346135696, time.UTC),
			content: []byte(""),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/circuit"].(os.FileInfo),
		fs["/package"].(os.FileInfo),
	}
	fs["/circuit"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/circuit/NOTE.txt"].(os.FileInfo),
		fs["/circuit/manifest.yaml"].(os.FileInfo),
		fs["/circuit/terraform.tf"].(os.FileInfo),
		fs["/circuit/values.hcl"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰FileInfo:
		return &vfsgen۰File{
			vfsgen۰FileInfo: f,
			Reader:          bytes.NewReader(f.content),
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰FileInfo is a static definition of an uncompressed file (because it's not worth gzip compressing).
type vfsgen۰FileInfo struct {
	name    string
	modTime time.Time
	content []byte
}

func (f *vfsgen۰FileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰FileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰FileInfo) NotWorthGzipCompressing() {}

func (f *vfsgen۰FileInfo) Name() string       { return f.name }
func (f *vfsgen۰FileInfo) Size() int64        { return int64(len(f.content)) }
func (f *vfsgen۰FileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰FileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰FileInfo) IsDir() bool        { return false }
func (f *vfsgen۰FileInfo) Sys() interface{}   { return nil }

// vfsgen۰File is an opened file instance.
type vfsgen۰File struct {
	*vfsgen۰FileInfo
	*bytes.Reader
}

func (f *vfsgen۰File) Close() error {
	return nil
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
