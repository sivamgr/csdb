package main

import (
	"encoding/csv"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func createDirForFile(filepath string) {
	dir := path.Dir(filepath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func utc2fileid(t uint32) uint32 {
	return (t / filecap()) * filecap()
}

func utc2fileidnext(t uint32) uint32 {
	return utc2fileid(t) + filecap()
}

func filecap() uint32 {
	return 86400 * 5
}

func fileNameWithNoDirNoExt(fp string) string {
	f := filepath.Base(fp)
	return strings.TrimSuffix(f, filepath.Ext(f))
}

// exists returns whether the given file or directory exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func str2float64(f string) float64 {
	if s, err := strconv.ParseFloat(f, 64); err == nil {
		return s
	}
	return 0
}

func str2u32(f string) uint32 {
	if s, err := strconv.ParseUint(f, 10, 32); err == nil {
		return uint32(s)
	}
	return 0
}
