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

// TODO out of place
func execute(obj Exec) error {
	defer obj.close() // TODO must call in any way!!!

	if err := obj.run(); err != nil {
		return err
	}

	return nil
}

func run(objs ...Exec) error {
	var wg sync.WaitGroup

	// Phase 1. PREPARING
	for i, obj := range objs {
		err := obj.prepare()
		if err != nil {
			// undo previous and current
			for _, o := range objs[:i+1] {
				o.close()
			}
			return err
		}
		// final obj's healthcheck
		err = obj.check()
		if err != nil {
			// undo previous and current
			for _, o := range objs[:i+1] {
				o.close()
			}
			return err
		}

		wg.Add(1)
	}

	// Phase 2. running
	// Run objects (layers) in order
	// TODO make new runner interface and add
	// .run for internal running
	// .Run for public run
	// .execute - ??? something with separate goroutine, err channels and context logic ???
	errch := make(chan error)

	for _, obj := range objs {

		go func(obj Exec) {
			defer wg.Done()

			errch <- execute(obj)
		}(obj)

	}

	select {
	case err := <-errch:
		if err != nil {
			// TODO cancel from context
			return err
		}
	}

	wg.Wait()

	return nil
}






// https://play.golang.org/p/SEXBheyHnt6
// https://play.golang.org/p/Zy7BpvwLlqg
// func parallel(fs ...func() error) error {
// }
