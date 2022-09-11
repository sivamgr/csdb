package main

import "fmt"

func readCSD0(fp string) []cs {
	c := new([]cs)
	readMsgpackFile(fp, c)
	return *c
}

func fetchSymbolData(sym string, ft uint32, tt uint32) []cs {
	c := make([]cs, 0)

	for t := utc2fileid(ft); t <= utc2fileid(tt); t = utc2fileidnext(t) {
		fp := fmt.Sprintf("%s/%s/%d.csd0", dataPath, sym, t)
		if exists(fp) {
			tc := readCSD0(fp)
			for _, it := range tc {
				if (it.UTC >= ft) && (it.UTC <= tt) {
					c = append(c, it)
				}
			}
		}
	}
	return c
}
