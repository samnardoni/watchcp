package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"time"
)

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

func instantTicker() chan time.Time {
	c := make(chan time.Time)
	go func() {
		c <- time.Now()
		t := time.NewTicker(1 * time.Second)
		for x := range t.C {
			c <- x
		}
	}()
	return c
}

func waitStat(file string) time.Time {
	for _ = range instantTicker() {

		stat, err := os.Stat(file)
		if err != nil {
			continue
		}

		return stat.ModTime()
	}

	return time.Time{}
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		println("Usage: watchcp src dst")
		return
	}

	src := flag.Arg(0)
	dst := flag.Arg(1)

	for {
		// Wait until src file becomes available
		srcStat := waitStat(src)

		// Wait until dst directory becomes available
		waitStat(filepath.Dir(dst))

		// Copy if dst file does not exist or src is newer then dst
		dstStat, err := os.Stat(dst)

		shouldCopy := err != nil || dstStat.ModTime().Before(srcStat)

		if !shouldCopy {
			continue
		}

		println("Copying", src, "to", dst, "...")

		_, err = copyFile(src, dst)
		if err != nil {
			println("Error", err.Error())
		}
	}
}
