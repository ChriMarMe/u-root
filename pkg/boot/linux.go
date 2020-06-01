// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/uio"
)

// LinuxImage implements OSImage for a Linux kernel + initramfs.
type LinuxImage struct {
	Name string

	Kernel  io.ReaderAt
	Initrd  io.ReaderAt
	Cmdline string
}

func module(r io.ReaderAt) map[string]interface{} {
	m := make(map[string]interface{})
	if f, ok := r.(curl.File); ok {
		m["url"] = f.URL().String()
	} else if f, ok := r.(fmt.Stringer); ok {
		m["stringer"] = f.String()
	}
	return m
}

// JSONMap is implemented only in order to compare LinuxImages in tests.
//
// It should be json-encodable and decodable.
func (li *LinuxImage) JSONMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["image_type"] = "linux"
	m["name"] = li.Name
	m["cmdline"] = li.Cmdline
	if li.Kernel != nil {
		m["kernel"] = module(li.Kernel)
	}
	if li.Initrd != nil {
		m["initrd"] = module(li.Initrd)
	}
	return m
}

func (li *LinuxImage) MarshalJSON() ([]byte, error) {
	return json.Marshal(li.JSONMap())
}

var _ OSImage = &LinuxImage{}

// Label returns either the Name or a short description.
func (li *LinuxImage) Label() string {
	if len(li.Name) > 0 {
		return li.Name
	}
	return fmt.Sprintf("Linux(kernel=%s, initrd=%s)", li.Kernel, li.Initrd)
}

// String prints a human-readable version of this linux image.
func (li *LinuxImage) String() string {
	return fmt.Sprintf("LinuxImage(\n  Name: %s\n  Kernel: %s\n  Initrd: %s\n  Cmdline: %s\n)\n", li.Name, li.Kernel, li.Initrd, li.Cmdline)
}

func copyToFile(r io.Reader) (*os.File, error) {
	f, err := ioutil.TempFile("", "nerf-netboot")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return nil, err
	}
	if err := f.Sync(); err != nil {
		return nil, err
	}

	readOnlyF, err := os.Open(f.Name())
	if err != nil {
		return nil, err
	}
	return readOnlyF, nil
}

// Load implements OSImage.Load and kexec_load's the kernel with its initramfs.
func (li *LinuxImage) Load(verbose bool) error {
	if li.Kernel == nil {
		return errors.New("LinuxImage.Kernel must be non-nil")
	}

	kernel, initrd := uio.Reader(li.Kernel), uio.Reader(li.Initrd)
	if verbose {
		// In verbose mode, print a dot every 5MiB. It is not pretty,
		// but it at least proves the files are still downloading.
		progress := func(r io.Reader) io.Reader {
			return &uio.ProgressReader{
				R:        r,
				Symbol:   ".",
				Interval: 5 * 1024 * 1024,
				W:        os.Stdout,
			}
		}
		kernel = progress(kernel)
		initrd = progress(initrd)
	}

	// It seams inefficient to always copy, in particular when the reader
	// is an io.File but that's not sufficient, os.File could be a socket,
	// a pipe or some other strange thing. Also kexec_file_load will fail
	// (similar to execve) if anything as the file opened for writing.
	// That's unfortunately something we can't guarantee here - unless we
	// make a copy of the file and dump it somewhere.
	k, err := copyToFile(kernel)
	if err != nil {
		return err
	}
	defer k.Close()

	var i *os.File
	if li.Initrd != nil {
		i, err = copyToFile(initrd)
		if err != nil {
			return err
		}
		defer i.Close()
	}

	log.Printf("Kernel: %s", k.Name())
	if i != nil {
		log.Printf("Initrd: %s", i.Name())
	}
	log.Printf("Command line: %s", li.Cmdline)
	return kexec.FileLoad(k, i, li.Cmdline)
}
