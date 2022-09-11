package main

import (
	"fmt"
	"sort"
)

type WrBucket struct {
	s  string
	c  []cs
	t0 uint32
	tn uint32
}

func (w *WrBucket) Add(c cs) bool {
	n := len(w.c)
	if n == 0 {
		w.c = append(w.c, c)
		w.t0 = utc2fileid(c.UTC)
		w.tn = utc2fileidnext(w.t0)
		return true
	}

	if (c.UTC < w.t0) || (c.UTC >= w.tn) {
		return false
	}

	if w.c[n-1].UTC >= c.UTC {
		return true
	}

	w.c = append(w.c, c)
	return true
}

func (w *WrBucket) Flushwrite() {
	if len(w.c) == 0 {
		return
	}

	fp := fmt.Sprintf("%s/%s/%d.csd0", dataPath, w.s, w.t0)
	createDirForFile(fp)
	writeMsgpackFile(fp, w.c)
	w.c = w.c[:0]
}

type SymDataWriter struct {
	s string
	c []cs
}

func (w *SymDataWriter) Add(n cs) {
	w.c = append(w.c, n)
}

func (w *SymDataWriter) Flushwrite() {
	if len(w.c) == 0 {
		return
	}

	tmin := w.c[0].UTC
	tmax := w.c[0].UTC

	for _, it := range w.c {
		if it.UTC > tmax {
			tmax = it.UTC
		}
		if it.UTC < tmin {
			tmin = it.UTC
		}
	}

	for t := utc2fileid(tmin); t <= utc2fileid(tmax); t = utc2fileidnext(t) {
		fp := fmt.Sprintf("%s/%s/%d.csd0", dataPath, w.s, t)
		if exists(fp) {
			cs := readCSD0(fp)
			if len(cs) > 0 {
				w.c = append(w.c, cs...)
			}
		}
	}

	sort.Slice(w.c, func(i, j int) bool {
		return w.c[i].UTC < w.c[j].UTC
	})

	wr := WrBucket{s: w.s}
	for _, c := range w.c {
		if !wr.Add(c) {
			wr.Flushwrite()
			wr.Add(c)
		}
	}
	wr.Flushwrite()
	w.c = w.c[:0]
}

type DBWriter struct {
	m map[string]*SymDataWriter
}

func (w *DBWriter) Add(s string, c cs) {
	if w.m == nil {
		w.m = make(map[string]*SymDataWriter)
	}
	if _, ok := w.m[s]; !ok {
		w.m[s] = new(SymDataWriter)
		w.m[s].s = s
	}

	w.m[s].Add(c)
}

func (w *DBWriter) Flushwrite() {
	for _, w := range w.m {
		w.Flushwrite()
	}
}
