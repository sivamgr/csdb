package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func arr2cs(v []string) cs {
	ts, _ := time.Parse("20060102 15:04 MST", v[0]+" "+v[1]+" IST")
	return cs{
		UTC: uint32(ts.Unix()),
		O:   uint32(math.Round(str2float64(v[2]) * 100)),
		H:   uint32(math.Round(str2float64(v[3]) * 100)),
		L:   uint32(math.Round(str2float64(v[4]) * 100)),
		C:   uint32(math.Round(str2float64(v[5]) * 100)),
		V:   str2u32(v[6]),
		OI:  str2u32(v[7]),
	}
}

func processInTextFile(fp string) {
	dbWrite := DBWriter{}

	log.Printf("Processing Text, %s\n", fp)
	s := fileNameWithNoDirNoExt(fp)
	d := readCsvFile(fp)
	for _, r := range d {
		fmt.Printf("%d ", len(r))
		if len(r) == 9 {
			cs := arr2cs(r[1:])
			s = r[0]

			fmt.Printf("%v\n", cs)
			dbWrite.Add(s, cs)

		} else if len(r) == 8 {
			cs := arr2cs(r)
			fmt.Printf("%v\n", cs)
			dbWrite.Add(s, cs)
		}
	}
	dbWrite.Flushwrite()
	os.Remove(fp)
}

// Closure to address file descriptors issue with all the deferred .Close() methods
func extractAndProcess(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	fp := filepath.Join(dest, f.Name)

	if _, err := os.Stat(filepath.Dir(fp)); os.IsNotExist(err) {
		log.Println("creating dir, ", filepath.Dir(fp))
		os.MkdirAll(filepath.Dir(fp), os.ModePerm)
	}

	// Check for ZipSlip (Directory traversal)
	if !strings.HasPrefix(fp, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", fp)
	}

	if f.FileInfo().IsDir() {
		os.MkdirAll(fp, f.Mode())
	} else {
		os.MkdirAll(filepath.Dir(fp), f.Mode())
		f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
		//log.Println("saved as, ", fp)
		processInFile(fp)
	}
	return nil
}

func unzipAndProcess(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	for _, f := range r.File {
		//log.Println("Unzipping", f.Name)
		err := extractAndProcess(f, dest)
		if err != nil {
			log.Printf("Error Unzipping, %v\n", err)
			return err
		}

	}

	return nil
}

func processInZipFile(fp string) {
	dir, _ := os.MkdirTemp("", "*")
	log.Printf("Processing zip, %s. Extracting into temp dir, %s\n", fp, dir)
	unzipAndProcess(fp, dir)
	os.RemoveAll(dir)
	os.Remove(fp)

}

func processInFile(fp string) {
	if fi, err := os.Stat(fp); err == nil {
		if !fi.IsDir() {
			fpExt := path.Ext(fp)

			if fpExt == ".zip" {
				processInZipFile(fp)
			} else if (fpExt == ".txt") || (fpExt == ".csv") {
				processInTextFile(fp)
			} else {
				os.Remove(fp)
			}
		}
	}

}

func processIncomingFiles() {
	filepath.Walk(inboxPath, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		processInFile(fp)
		return nil
	})
}
