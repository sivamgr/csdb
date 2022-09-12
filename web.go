package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func inboxPutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only POST is supported in inbox")
		return
	}

	// Maximum upload of 10 MB files
	r.ParseMultipartForm(5 << 20)
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Error")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fp := fmt.Sprintf("%s%s", inboxPath, filepath.Base(handler.Filename))

	// Create file
	dst, err := os.Create(fp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		dst.Close()
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		dst.Close()
		return
	}
	dst.Close()
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	go processInFile(fp)
}

func symbolsListGetHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Symbols!")

	files, err := os.ReadDir(dataPath)
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	syms := make([]string, 0)
	for _, f := range files {
		fn := f.Name()
		if len(fn) > 1 {
			syms = append(syms, fn)
		}

	}
	json.NewEncoder(w).Encode(syms)
}

type candleItem struct {
	T  uint32  `json:"t"`
	O  float64 `json:"o"`
	H  float64 `json:"h"`
	L  float64 `json:"l"`
	C  float64 `json:"c"`
	V  uint32  `json:"v"`
	OI uint32  `json:"oi"`
}

type dataItem struct {
	S string       `json:"s"`
	C []candleItem `json:"c"`
}

func cs2ci(c cs) candleItem {
	return candleItem{
		T:  c.UTC,
		O:  float64(c.O) / 100.0,
		H:  float64(c.H) / 100.0,
		L:  float64(c.L) / 100.0,
		C:  float64(c.C) / 100.0,
		V:  c.V,
		OI: c.OI,
	}
}

func dataGetHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Data!")
	//fmt.Printf("%v\n", r.URL.Path)
	values := r.URL.Query()
	syms := values["s"]
	f := values["f"]
	t := values["t"]

	tt := uint32(time.Now().Unix())
	ft := tt - (86400 * 7)

	if len(syms) == 0 {
		return
	}

	if len(f) > 0 {
		v := str2u32(f[0])
		if v > 0 {
			if (v>19800101) && (v<20390101) {
				ds := fmt.Sprintf("%04d-%02d-%02dT00:00:00.000Z",v/10000,(v/100)%100,v%100)
				ts,	 err := time.Parse("2006-01-02T15:04:05.000Z", ds)
				if err == nil {
					v = uint32(ts.Unix())
				}
			}
			ft = v
		}
	}
	if len(t) > 0 {
		v := str2u32(f[0])
		if v > 0 {
			if (v>19800101) && (v<20390101) {
				ds := fmt.Sprintf("%04d-%02d-%02dT23:59:59.000Z",v/10000,(v/100)%100,v%100)
				ts,	 err := time.Parse("2006-01-02T15:04:05.000Z", ds)
				if err == nil {
					v = uint32(ts.Unix())
				}
			}
			tt = v
		}
	}

	resp := make([]dataItem, 0)
	for _, sym := range syms {
		c := fetchSymbolData(sym, ft, tt)
		candles := make([]candleItem, 0)
		for _, csit := range c {
			candles = append(candles, cs2ci(csit))
		}
		resp = append(resp, dataItem{S: sym, C: candles})
	}
	/*
		jData, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "error")
			return
		}

		w.Write(jData)
	*/
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func startWebServer() {
	http.HandleFunc("/inbox", inboxPutHandler)
	http.HandleFunc("/symbols", symbolsListGetHandler)
	http.HandleFunc("/data", dataGetHandler)

	addr := os.Getenv("CSDB_SERVER_ADDRESS")
	if len(addr) == 0 {
		addr = ":8087"
	}
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
