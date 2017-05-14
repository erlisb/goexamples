package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func saveFile(fileName string, buf []byte) error {
	writer, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer writer.Close()

	if _, err = io.Copy(writer, bytes.NewReader(buf)); err != nil {
		return err
	}

	return nil
}

func decompress(fileName string) ([]byte, error) {
	// open the file
	fp, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	defer fp.Close()

	//
	reader, err := zlib.NewReader(fp)

	if err != nil {
		return nil, err
	}

	defer reader.Close()

	//
	buf, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func main() {
	flag.Parse() // get the arguments from command line

	fileName := flag.Arg(0)

	if fileName == "" {
		fmt.Println("Please sepcify a filename")
		os.Exit(1)
	}

	// decompress
	buf, err := decompress(fileName)

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		os.Exit(1)
	}

	// save to file
	if err := saveFile(path.Base(fileName)+".txt", buf); err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		os.Exit(1)
	}

	// output to screen
	fmt.Printf("%q\n", buf)

}
