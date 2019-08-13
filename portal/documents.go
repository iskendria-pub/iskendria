package main

import (
	"errors"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"io/ioutil"
	"os"
	"sync"
)

func initDocuments() {
	theDocuments = &documents{
		path: "./documents",
	}
	theDocuments.init()
}

type documents struct {
	mux  sync.Mutex
	path string
}

var theDocuments *documents

func (d *documents) init() {
	d.mux.Lock()
	defer d.mux.Unlock()
	_, err := os.Stat(d.path)
	switch {
	case err == nil:
		// Do nothing
	case os.IsNotExist(err):
		err := os.Mkdir(d.path, 0744)
		if err != nil {
			panic(err)
		}
	default:
		panic(err)
	}
}

func (d *documents) searchDescription(theHash string) ([]byte, bool, error) {
	if theHash == "" {
		return []byte{}, true, nil
	}
	d.mux.Lock()
	defer d.mux.Unlock()
	_, err := os.Stat(d.getPath(theHash))
	switch {
	case err == nil:
		result, readErr := d.readAndCheck(theHash)
		return result, true, readErr
	case os.IsNotExist(err):
		return []byte{}, false, nil
	default:
		return []byte{}, false, err
	}
}

func (d *documents) getPath(theHash string) string {
	return d.path + "/" + theHash
}

func (d *documents) readAndCheck(theHash string) ([]byte, error) {
	fname := d.getPath(theHash)
	resultBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return []byte{}, err
	}
	if model.HashBytes(resultBytes) != theHash {
		err := os.Remove(fname)
		if err != nil {
			panic(err)
		}
		return []byte{}, errors.New("Removed document because of hash mismatch, hash: " + theHash)
	}
	return resultBytes, nil
}

func (d *documents) save(theHash string, data []byte) error {
	if model.HashBytes(data) != theHash {
		return errors.New("Uploaded file does not have hash " + theHash)
	}
	return ioutil.WriteFile(d.getPath(theHash), data, 0644)
}

func (d *documents) open(theHash string) (*os.File, error) {
	return os.Open(d.getPath(theHash))
}
