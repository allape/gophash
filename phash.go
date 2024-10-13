package gophash

import (
	"bytes"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"unsafe"
)

/*
#cgo LDFLAGS: -L/usr/local/lib -lphash -lstdc++
#cgo CXXFLAGS: -std=c++11

#include <stdint.h>
#include <stdlib.h>

extern double ph_hammingdistance2(uint8_t *hashA, int lenA, uint8_t *hashB, int lenB);
extern uint8_t *ph_mh_imagehash(const char *filename, int *N, float alpha, float lvl);
*/
import "C"

type PHash struct {
	hash.Hash

	Logger   *log.Logger
	Filename string

	Alpha float32
	Level float32
}

func (d *PHash) Size() int {
	return 72 // ulong64
}

func (d *PHash) BlockSize() int {
	return d.Size()
}

func (d *PHash) Reset() {}

func (d *PHash) Write(_ []byte) (int, error) {
	return 0, nil
}

func (d *PHash) Sum(b []byte) []byte {
	defer func() {
		_ = recover()
	}()
	return d.sum(b)
}

func (d *PHash) sum(b []byte) []byte {
	filename := d.Filename
	if b != nil {
		tmp, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("phash-*-%s", d.Filename))
		if err != nil {
			d.Logger.Println("Unable to create temporary file:", err)
			return nil
		}
		defer func() {
			_ = tmp.Close()
			_ = os.Remove(tmp.Name())
		}()
		_, err = io.Copy(tmp, bytes.NewReader(b))
		if err != nil {
			d.Logger.Println("Unable to write temporary file:", err)
			return nil
		}
		filename = tmp.Name()
	} else {
		stat, err := os.Stat(filename)
		if err != nil {
			d.Logger.Println("Unable to stat file:", err)
			return nil
		} else if stat.IsDir() {
			d.Logger.Println(filename, "is a directory")
			return nil
		}
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cN := C.int(0)

	res := C.ph_mh_imagehash(cFilename, &cN, C.float(d.Alpha), C.float(d.Level))
	defer C.free(unsafe.Pointer(res))

	hashBytes := C.GoBytes(unsafe.Pointer(res), cN)

	return hashBytes
}

func New(filename string, logger *log.Logger) *PHash {
	if logger == nil {
		logger = log.New(io.Discard, "", log.LstdFlags)
	}
	return &PHash{
		Logger:   logger,
		Filename: filename,
	}
}
