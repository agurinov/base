package pipeline

import (
	"io"
	"sync"
)

// piping establishes pipe connections between IO processes (Layer)
// the first layer accepts as stdin the input buffer
// the last layer puts into output buffer his stdout
func piping(input io.ReadCloser, output io.WriteCloser, layers ...Able) error {
	// main logic that create pairs of (io.ReadCloser, io.WriteCloser)
	// but with offset to another layer
	// for example
	// layer 1: (input, io.WriteCloser 1)
	// layer 2: (io.ReadCloser 1, io.WriteCloser 2)
	// layer 3: (io.ReadCloser 2, io.WriteCloser 3)
	// layer 4: (io.ReadCloser 3, output)
	// and call each layer's .pipe() method
	for i := 0; i < len(layers); i++ {
		if i == 0 {
			// case this layer first
			layers[i].setStdin(input)
		}
		if i == len(layers)-1 {
			// case this layer last
			layers[i].setStdout(output)
		} else {
			// this is intermediate layer, need piping
			r, w := io.Pipe()
			layers[i].setStdout(w)
			layers[i+1].setStdin(r)
		}
	}

	return nil
}

func run(objs ...Exec) error {
	var wg sync.WaitGroup

	// TODO join this 2 cycles with one and defer!
	// TODO test with fails
	// prepare all objs (prepare hook)
	for _, obj := range objs {
		if err := obj.prepare(); err != nil {
			return err
		}

		wg.Add(1)
	}

	// Run pipeline
	for _, obj := range objs {
		go func(obj Exec) {
			defer obj.close()
			defer wg.Done()

			obj.run()
		}(obj)
	}

	wg.Wait()

	return nil
}
