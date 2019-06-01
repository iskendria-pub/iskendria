package main

import (
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

func (d *documents) searchDescription(theHash string) (string, bool, error) {
	d.mux.Lock()
	defer d.mux.Unlock()
	_, err := os.Stat(d.getPath(theHash))
	switch {
	case err == nil:
		result, readErr := readDocumentFile(d.getPath(theHash))
		return result, true, readErr
	case os.IsNotExist(err):
		return "", false, nil
	default:
		return "", false, err
	}
}

func (d *documents) getPath(theHash string) string {
	return d.path + "/" + theHash
}

func readDocumentFile(fname string) (string, error) {
	resultBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return string(resultBytes), nil
}

func (d *documents) Save(theHash string, data []byte) error {
	return ioutil.WriteFile(d.getPath(theHash), data, 0644)
}
