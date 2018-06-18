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

func execute(obj Exec) error {
	defer obj.close()

	if err := obj.run(); err != nil {
		return err
	}

	return nil
}

func prepare(objs ...Exec) (err error) {
	var i int

	// backwards closing and resetting if error exists
	defer func() {
		if err != nil {
			// need to backwards
			for ; i >= 0; i-- {
				objs[i].close()
			}
		}
	}()

	// iterate over layers
	for ; i < len(objs); i++ {
		// try to prepare obj
		if err = objs[i].prepare(); err != nil {
			return
		}
		// final obj's healthcheck
		if err = objs[i].check(); err != nil {
			return
		}
	}

	return
}

func run(objs ...Exec) error {
	// Phase 1. PREPARING
	// in case of error it will be rolled back to initial incoming state
	if err := prepare(objs...); err != nil {
		return err
	}

	// Phase 2. running. Here ALL objs ready and checked
	var wg sync.WaitGroup
	wg.Add(len(objs))
	ch := make(chan error)

	for _, obj := range objs {

		go func(obj Exec) {
			defer wg.Done()

			ch <- execute(obj)
		}(obj)

	}

	select {
	case err := <-ch:
		if err != nil {
			return err
		}
	}

	wg.Wait()

	return nil
}





// https://play.golang.org/p/Djv52XGnbur
// https://play.golang.org/p/SEXBheyHnt6
// https://play.golang.org/p/Zy7BpvwLlqg
// func parallel(fs ...func() error) error {
// }
