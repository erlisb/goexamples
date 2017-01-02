package main

import (
	"bytes"
	"fmt"
	"testing"
)

func BenchmarkGobGzipV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		person := GetPerson()
		person2 := Person{}

		if b, err := EncodeGobThenGzip(&person); err != nil {
			fmt.Printf("err: %v\n\n", err)
		} else if err := UngzipThenDecodeGob(bytes.NewReader(b), &person2); err != nil {
			fmt.Printf("err: %v\n\n", err)
		}
	}
}
