package goconf

import "os"

func pathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
