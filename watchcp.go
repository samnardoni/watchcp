package main

import (
	"flag"
	"io"
	"os"
	"time"
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		println("Usage: watchcp src dst")
		return
	}

	src := flag.Arg(0)
	dst := flag.Arg(1)

	for {
		copyFileIfNewer(src, dst)
		time.Sleep(1 * time.Second)
	}
}

func copyFileIfNewer(src, dst string) {
	if shouldCopy(src, dst) {
		println("Copying", src, "to", dst, "...")

		_, err := copyFile(src, dst)
		if err != nil {
			println("Error copying file:", err.Error())
		}
	}
}

func shouldCopy(src, dst string) bool {
	srcStat, err := os.Stat(src)
	if err != nil {
		return false
	}

	dstStat, err := os.Stat(dst)
	if err != nil {
		return true
	}

	return srcStat.ModTime().After(dstStat.ModTime())
}

func copyFile(src, dst string) (int64, error) {
	sf, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	defer sf.Close()
	if err != nil {
		return 0, err
	}

	df, err := os.OpenFile(dst, os.O_TRUNC|os.O_CREATE, os.ModePerm)
	defer df.Close()
	if err != nil {
		return 0, err
	}

	return io.Copy(df, sf)
}
