package pipeline

import (
	"io"
	"sync"
)

// piping establishes pipe connections between IO processes (Able)
// the first obj accepts as stdin the input buffer
// the last obj puts into output buffer his stdout
func piping(input io.ReadCloser, output io.WriteCloser, objs ...Able) error {
	// main logic that create pairs of (io.ReadCloser, io.WriteCloser)
	// but with offset to another obj
	// for example
	// obj 1: (input, io.WriteCloser 1)
	// obj 2: (io.ReadCloser 1, io.WriteCloser 2)
	// obj 3: (io.ReadCloser 2, io.WriteCloser 3)
	// obj 4: (io.ReadCloser 3, output)
	for i := 0; i < len(objs); i++ {
		if i == 0 {
			// case this obj first
			objs[i].setStdin(input)
		}
		if i == len(objs)-1 {
			// case this obj last
			objs[i].setStdout(output)
		} else {
			// this is intermediate obj, need piping
			r, w := io.Pipe()
			objs[i].setStdout(w)
			objs[i+1].setStdin(r)
		}
	}

	return nil
}

func run(objs ...Exec) error {
	var wg sync.WaitGroup

	// TODO join this 2 cycles with one and defer!
	// TODO test with fails
	for _, obj := range objs {
		// prepare obj
		if err := obj.prepare(); err != nil {
			return err
		}

		// final obj's healthcheck
		if err := obj.check(); err != nil {
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
