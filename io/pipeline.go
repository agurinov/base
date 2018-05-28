package io

import (
	"io"
	// "io/ioutil"
	"fmt"
	"log"
	"os/exec"
	// "bytes"
)

type Pipeline []io.ReadWriter


type Layer struct {
	cmd *exec.Cmd
}
func (l *Layer) Write(p []byte) (n int, err error) {
	return l.cmd.Stdout.Write(p)
}
func (l *Layer) Read(p []byte) (n int, err error) {
	return l.cmd.Stdin.Read(p)
}


func chaining(input io.Reader, output io.Writer, layer *exec.Cmd) {
	// buf, err := ioutil.ReadAll(input)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	buf := make([]byte, 5)

	reader, writer := io.Pipe()

	layer.Stdin = reader
	layer.Stdout = writer

	layer.Start()

	io.Copy(writer, input)

	if err := layer.Wait(); err != nil {
		log.Fatal(err)
	}

	writer.Close()

	// final response
	io.Copy(output, reader)
	reader.Close()
}



// func main() {
//     c1 := exec.Command("ls")
//     c2 := exec.Command("wc", "-l")
//
//     r, w := io.Pipe()
//     c1.Stdout = w
//     c2.Stdin = r
//
//     var b2 bytes.Buffer
//     c2.Stdout = &b2
//
//     c1.Start()
//     c2.Start()
//     c1.Wait()
//     w.Close()
//     c2.Wait()
//     io.Copy(os.Stdout, &b2)
// }
