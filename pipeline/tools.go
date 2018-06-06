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

func run(layers ...Exec) error {
	var wg sync.WaitGroup

	// TODO join this 2 cycles with one and defer!
	// TODO test with fails
	// prepare all layers (prepare hook)
	for _, layer := range layers {
		if err := layer.prepare(); err != nil {
			return err
		}

		wg.Add(1)
	}

	// Run pipeline
	for _, layer := range layers {
		go func(layer Exec) {
			defer layer.Close()
			defer wg.Done()

			layer.Run()
		}(layer)
	}

	wg.Wait()

	return nil
}

// // start invokes the layer's .Start() method in the order of the queue
// func start(layers []*exec.Cmd) (err error) {
// 	// start the pipeline
// 	for _, layer := range layers {
// 		if err := layer.Start(); err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// // run causes processes in turn
// // waiting for the previous stdout layer and picks it into the next one
// func run(layers []*exec.Cmd, pipeWriters []*io.PipeWriter) (err error) {
// 	// original idea:
// 	// https://gist.github.com/tyndyll/89fbb2c2273f83a074dc
// 	for i, layer := range layers {
// 		if err := layer.Wait(); err != nil {
// 			return err
// 		}
//
// 		// if next layers in queue exists -> close pipe
// 		if i < len(layers)-1 {
// 			if err := pipeWriters[i].Close(); err != nil {
// 				return err
// 			}
// 		}
// 	}
//
// 	return nil
// }
