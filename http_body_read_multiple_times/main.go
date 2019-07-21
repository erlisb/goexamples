// Reference:
// https://medium.com/@xoen/golang-read-from-an-io-readwriter-without-loosing-its-content-2c6911805361
// https://github.com/gin-gonic/gin/issues/1295
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Person struct {
	Name string
}

func srvExample1(w http.ResponseWriter, r *http.Request) {
	// Read the content
	//var bodyBytes []byte

	bodyBytes := make([]byte, 1024)
	numOfBytes := 0
	var err error

	if r.Body != nil {
		// http.Response.Body is of type io.ReadCloser, which can only be read once.
		// When you read from io.ReadCloser it drains it.
		// Once you read from it, the content is gone. You can’t read from it a second time.

		if numOfBytes, err = r.Body.Read(bodyBytes); err != nil {
			if err != io.EOF {
				fmt.Fprintf(w, "Error: %s!", err.Error())
				return
			}
		}
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes[0:numOfBytes]))

	// Read the r.Body one more time
	p1 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p1); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	fmt.Fprintf(w, "%#v\n", p1)
	//w.Write([]byte("Hello"))
}

func srvExample2(w http.ResponseWriter, r *http.Request) {
	var err error
	var bodyBytes []byte

	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Read the r.Body one more time
	p1 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p1); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	fmt.Fprintf(w, "%#v\n", p1)
}
func srvExample3(w http.ResponseWriter, r *http.Request) {
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))

	// TeeReader returns a Reader that writes to w what it reads from r.
	reader := io.TeeReader(r.Body, buf)

	//
	p1 := Person{}

	// Note: using reader instead of r.Body
	if err := json.NewDecoder(reader).Decode(&p1); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(buf)

	// Read the r.Body one more time
	p2 := Person{}

	if err = json.NewDecoder(r.Body).Decode(&p2); err != nil {
		fmt.Fprintf(w, "Error: %s!", err.Error())
		return
	}

	fmt.Fprintf(w, "%#v\n%#v\n", p1, p2)
}

func main() {
	// Example1 uses r.Body.Read()
	// curl http://localhost:8080/Example1 --data '{"Name":"asdf2"}'
	http.HandleFunc("/Example1", srvExample1)

	// Example2 uses ioutil.ReadAll()
	// curl http://localhost:8080/Example2 --data '{"Name":"asdf2"}'
	http.HandleFunc("/Example2", srvExample2)

	// Example3 uses io.TeeReader(()
	http.HandleFunc("/Example3", srvExample3)

	http.ListenAndServe(":8080", nil)
}
