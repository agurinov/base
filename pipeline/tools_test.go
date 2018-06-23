package pipeline

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

func checkObjsByMatrixState(t *testing.T, matrix [][]int, objs []interface{}) {
	t.Log(len(matrix), len(objs))
	// check sizes
	if len(matrix) != len(objs) {
		t.Fatalf("CheckByMatrix: inappropriate sizes, len(matrix)==%d != len(objs)==%d", len(matrix), len(objs))
	}

	// table tests
	for i, obj := range objs {
		if obj.(*execObj).countPrepare != matrix[i][0] {
			t.Fatalf("objs[%d].countPrepare: expected \"%d\", got \"%d\"", i, matrix[i][0], obj.(*execObj).countPrepare)
		}
		if obj.(*execObj).countCheck != matrix[i][1] {
			t.Fatalf("objs[%d].countCheck: expected \"%d\", got \"%d\"", i, matrix[i][1], obj.(*execObj).countCheck)
		}
		if obj.(*execObj).countRun != matrix[i][2] {
			t.Fatalf("objs[%d].countRun: expected \"%d\", got \"%d\"", i, matrix[i][2], obj.(*execObj).countRun)
		}
		if obj.(*execObj).countClose != matrix[i][3] {
			t.Fatalf("objs[%d].countClose: expected \"%d\", got \"%d\"", i, matrix[i][3], obj.(*execObj).countClose)
		}
	}
}

// execObj is special `fake` Layer
type execObj struct {
	countPrepare int
	countCheck   int
	countRun     int
	countClose   int

	mockFailPrepare bool
	mockFailCheck   bool
	mockFailRun     bool
	mockFailClose   bool
}

func (o *execObj) prepare() error {
	o.countPrepare++

	if o.mockFailPrepare {
		return errors.New("prepare failed")
	}

	return nil
}
func (o *execObj) check() error {
	o.countCheck++

	if o.mockFailCheck {
		return errors.New("check failed")
	}

	return nil
}
func (o *execObj) run() error {
	o.countRun++

	if o.mockFailRun {
		return errors.New("run failed")
	}

	return nil
}
func (o *execObj) close() error {
	o.countClose++

	if o.mockFailClose {
		return errors.New("close failed")
	}

	return nil
}
func (o *execObj) setStdin(reader io.ReadCloser) {
}
func (o *execObj) setStdout(writer io.WriteCloser) {
}

func TestToCloser(t *testing.T) {
	noCloser := bytes.NewBuffer([]byte{})
	nativeReadCloser := toReadCloser(noCloser)
	nativeWriteCloser := toWriteCloser(noCloser)

	t.Run("Read", func(t *testing.T) {
		t.Run("native", func(t *testing.T) {
			oldPtr := reflect.ValueOf(nativeReadCloser).Pointer()
			newPtr := reflect.ValueOf(toReadCloser(nativeReadCloser)).Pointer()

			// native ReadCloser will return without any injections -> same pointer
			if oldPtr != newPtr {
				t.Error("Unexpected pointer")
			}
		})

		t.Run("obtained", func(t *testing.T) {
			oldPtr := reflect.ValueOf(noCloser).Pointer()
			newPtr := reflect.ValueOf(toReadCloser(noCloser)).Pointer()

			// no ReadCloser will be returned with injection of .close() method (just return nil)
			// -> different pointer
			if oldPtr == newPtr {
				t.Error("Unexpected pointer")
			}
		})
	})

	t.Run("Write", func(t *testing.T) {
		t.Run("native", func(t *testing.T) {
			oldPtr := reflect.ValueOf(nativeWriteCloser).Pointer()
			newPtr := reflect.ValueOf(toWriteCloser(nativeWriteCloser)).Pointer()

			// native WriteCloser will return without any injections -> same pointer
			if oldPtr != newPtr {
				t.Error("Unexpected pointer")
			}
		})

		t.Run("obtained", func(t *testing.T) {
			oldPtr := reflect.ValueOf(noCloser).Pointer()
			newPtr := reflect.ValueOf(toWriteCloser(noCloser)).Pointer()

			// no WriteCloser will be returned with injection of .close() method (just return nil)
			// -> different pointer
			if oldPtr == newPtr {
				t.Error("Unexpected pointer")
			}
		})
	})
}

func TestPiping(t *testing.T) {
	input := toReadCloser(bytes.NewBuffer([]byte{}))
	output := toWriteCloser(bytes.NewBuffer([]byte{}))
	inputPtr := reflect.ValueOf(input).Pointer()
	outputPtr := reflect.ValueOf(output).Pointer()

	t.Run("tcp", func(t *testing.T) {
		if inputPtr == outputPtr {
			t.Fatal("Unexpected same pointers for input and output")
		}

		t.Run("len==1", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewTCPSocket("example.com:80"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdinPtr := reflect.ValueOf(layers[0].(*tcp).stdin).Pointer()
			stdoutPtr := reflect.ValueOf(layers[0].(*tcp).stdout).Pointer()

			if stdinPtr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			if stdoutPtr != outputPtr {
				t.Fatal("layers[0]: unexpected stdout")
			}
		})

		t.Run("len==2", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewTCPSocket("example.com:80"),
				NewTCPSocket("domain.com:22"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdin1Ptr := reflect.ValueOf(layers[0].(*tcp).stdin).Pointer()
			// stdout1Ptr := reflect.ValueOf(layers[0].(*tcp).stdout).Pointer()
			// stdin2Ptr := reflect.ValueOf(layers[1].(*tcp).stdin).Pointer()
			stdout2Ptr := reflect.ValueOf(layers[1].(*tcp).stdout).Pointer()

			if stdin1Ptr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			// if stdout1Ptr != stdin2Ptr {
			// 	t.Fatal("layers[0]: unexpected stdout")
			// }
			//
			// if stdin2Ptr != stdout1Ptr {
			// 	t.Fatal("layers[1]: unexpected stdin")
			// }
			if stdout2Ptr != outputPtr {
				t.Fatal("layers[1]: unexpected stdout")
			}
		})
	})

	t.Run("process", func(t *testing.T) {
		if inputPtr == outputPtr {
			t.Fatal("Unexpected same pointers for input and output")
		}

		t.Run("len==1", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewProcess("pwd"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdinPtr := reflect.ValueOf(layers[0].(*process).stdin).Pointer()
			stdoutPtr := reflect.ValueOf(layers[0].(*process).stdout).Pointer()

			if stdinPtr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			if stdoutPtr != outputPtr {
				t.Fatal("layers[0]: unexpected stdout")
			}
		})

		t.Run("len==2", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewProcess("pwd"),
				NewProcess("rev"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdin1Ptr := reflect.ValueOf(layers[0].(*process).stdin).Pointer()
			// stdout1Ptr := reflect.ValueOf(layers[0].(*process).stdout).Pointer()
			// stdin2Ptr := reflect.ValueOf(layers[1].(*process).stdin).Pointer()
			stdout2Ptr := reflect.ValueOf(layers[1].(*process).stdout).Pointer()

			if stdin1Ptr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			// if stdout1Ptr != stdin2Ptr {
			// 	t.Fatal("Unexpected stdout for [0] layer")
			// }
			//
			// if stdin2Ptr != stdout1Ptr {
			// 	t.Fatal("Unexpected stdin for [1] layer")
			// }
			if stdout2Ptr != outputPtr {
				t.Fatal("layers[1]: unexpected stdout")
			}
		})
	})

	t.Run("mix", func(t *testing.T) {
		t.Run("len==3", func(t *testing.T) {
			// layers for piping
			layers := []Able{
				NewProcess("pwd"),
				NewTCPSocket("example.com:80"),
				NewProcess("rev"),
			}
			// check for errors
			if err := piping(input, output, layers...); err != nil {
				t.Fatal(err)
			}
			// check layers stdio
			stdin1Ptr := reflect.ValueOf(layers[0].(*process).stdin).Pointer()
			// stdout1Ptr := reflect.ValueOf(layers[0].(*process).stdout).Pointer()
			//
			// stdin2Ptr := reflect.ValueOf(layers[1].(*tcp).stdin).Pointer()
			// stdout2Ptr := reflect.ValueOf(layers[1].(*tcp).stdout).Pointer()
			//
			// stdin3Ptr := reflect.ValueOf(layers[2].(*process).stdin).Pointer()
			stdout3Ptr := reflect.ValueOf(layers[2].(*process).stdout).Pointer()

			if stdin1Ptr != inputPtr {
				t.Fatal("layers[0]: unexpected stdin")
			}
			// if stdout1Ptr != stdin2Ptr {
			// 	t.Fatal("Unexpected stdout for [0] layer")
			// }
			//
			// if stdin2Ptr != stdout1Ptr {
			// 	t.Fatal("Unexpected stdin for [1] layer")
			// }
			if stdout3Ptr != outputPtr {
				t.Fatal("layers[2]: unexpected stdout")
			}
		})
	})
}

func TestPrepare(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		layers := []Exec{
			&execObj{},
			&execObj{},
			&execObj{},
		}
		// digits - count of invokes functions
		// digits from left to right:
		// countPrepare, countCheck, countRun, countClose
		flags := [][]int{
			[]int{1, 1, 0, 0},
			[]int{1, 1, 0, 0},
			[]int{1, 1, 0, 0},
		}

		// common errors
		if err := prepare(layers...); err != nil {
			t.Fatal(err)
		}

		// table tests
		for i, obj := range layers {
			if obj.(*execObj).countPrepare != flags[i][0] {
				t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
			}
			if obj.(*execObj).countCheck != flags[i][1] {
				t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
			}
			if obj.(*execObj).countRun != flags[i][2] {
				t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
			}
			if obj.(*execObj).countClose != flags[i][3] {
				t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
			}
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Run("prepare", func(t *testing.T) {
			layers := []Exec{
				&execObj{},
				&execObj{mockFailPrepare: true}, // backwards from here. i == 1
				&execObj{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := [][]int{
				[]int{1, 1, 0, 1},
				[]int{1, 0, 0, 1},
				[]int{0, 0, 0, 0},
			}

			// common errors
			err := prepare(layers...)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != "prepare failed" {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			for i, obj := range layers {
				if obj.(*execObj).countPrepare != flags[i][0] {
					t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
				}
				if obj.(*execObj).countCheck != flags[i][1] {
					t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
				}
				if obj.(*execObj).countRun != flags[i][2] {
					t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
				}
				if obj.(*execObj).countClose != flags[i][3] {
					t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
				}
			}
		})

		t.Run("check", func(t *testing.T) {
			layers := []Exec{
				&execObj{},
				&execObj{mockFailCheck: true}, // backwards from here. i == 1
				&execObj{mockFailPrepare: true},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := [][]int{
				[]int{1, 1, 0, 1},
				[]int{1, 1, 0, 1},
				[]int{0, 0, 0, 0},
			}

			// common errors
			err := prepare(layers...)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != "check failed" {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			for i, obj := range layers {
				if obj.(*execObj).countPrepare != flags[i][0] {
					t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
				}
				if obj.(*execObj).countCheck != flags[i][1] {
					t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
				}
				if obj.(*execObj).countRun != flags[i][2] {
					t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
				}
				if obj.(*execObj).countClose != flags[i][3] {
					t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
				}
			}
		})
	})
}

func TestExecute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		obj := &execObj{}

		// digits - count of invokes functions
		// digits from left to right:
		// countPrepare, countCheck, countRun, countClose
		flags := []int{0, 0, 1, 1}

		// common errors
		if err := execute(obj); err != nil {
			t.Fatal(err)
		}

		if obj.countPrepare != flags[0] {
			t.Fatalf("obj.countPrepare: expected \"%d\", got \"%d\"", flags[0], obj.countPrepare)
		}
		if obj.countCheck != flags[1] {
			t.Fatalf("obj.countCheck: expected \"%d\", got \"%d\"", flags[1], obj.countCheck)
		}
		if obj.countRun != flags[2] {
			t.Fatalf("obj.countRun: expected \"%d\", got \"%d\"", flags[2], obj.countRun)
		}
		if obj.countClose != flags[3] {
			t.Fatalf("obj.countClose: expected \"%d\", got \"%d\"", flags[3], obj.countClose)
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Run("run", func(t *testing.T) {
			obj := &execObj{mockFailRun: true, mockFailClose: true}

			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := []int{0, 0, 1, 1}

			// common errors
			err := execute(obj)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != "run failed" {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			if obj.countPrepare != flags[0] {
				t.Fatalf("obj.countPrepare: expected \"%d\", got \"%d\"", flags[0], obj.countPrepare)
			}
			if obj.countCheck != flags[1] {
				t.Fatalf("obj.countCheck: expected \"%d\", got \"%d\"", flags[1], obj.countCheck)
			}
			if obj.countRun != flags[2] {
				t.Fatalf("obj.countRun: expected \"%d\", got \"%d\"", flags[2], obj.countRun)
			}
			if obj.countClose != flags[3] {
				t.Fatalf("obj.countClose: expected \"%d\", got \"%d\"", flags[3], obj.countClose)
			}
		})

		t.Run("close", func(t *testing.T) {
			obj := &execObj{mockFailClose: true}

			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := []int{0, 0, 1, 1}

			// common errors
			err := execute(obj)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != "close failed" {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			if obj.countPrepare != flags[0] {
				t.Fatalf("obj.countPrepare: expected \"%d\", got \"%d\"", flags[0], obj.countPrepare)
			}
			if obj.countCheck != flags[1] {
				t.Fatalf("obj.countCheck: expected \"%d\", got \"%d\"", flags[1], obj.countCheck)
			}
			if obj.countRun != flags[2] {
				t.Fatalf("obj.countRun: expected \"%d\", got \"%d\"", flags[2], obj.countRun)
			}
			if obj.countClose != flags[3] {
				t.Fatalf("obj.countClose: expected \"%d\", got \"%d\"", flags[3], obj.countClose)
			}
		})
	})
}

func TestRun(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("len==1", func(t *testing.T) {
			layers := []Exec{
				&execObj{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := [][]int{
				[]int{1, 1, 1, 1},
			}

			// common errors
			if err := run(layers...); err != nil {
				t.Fatal(err)
			}

			// table tests
			for i, obj := range layers {
				if obj.(*execObj).countPrepare != flags[i][0] {
					t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
				}
				if obj.(*execObj).countCheck != flags[i][1] {
					t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
				}
				if obj.(*execObj).countRun != flags[i][2] {
					t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
				}
				if obj.(*execObj).countClose != flags[i][3] {
					t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
				}
			}
		})

		t.Run("len==3", func(t *testing.T) {
			layers := []Exec{
				&execObj{},
				&execObj{},
				&execObj{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := [][]int{
				[]int{1, 1, 1, 1},
				[]int{1, 1, 1, 1},
				[]int{1, 1, 1, 1},
			}

			// common errors
			if err := run(layers...); err != nil {
				t.Fatal(err)
			}

			// table tests
			for i, obj := range layers {
				if obj.(*execObj).countPrepare != flags[i][0] {
					t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
				}
				if obj.(*execObj).countCheck != flags[i][1] {
					t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
				}
				if obj.(*execObj).countRun != flags[i][2] {
					t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
				}
				if obj.(*execObj).countClose != flags[i][3] {
					t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
				}
			}
		})
	})

	t.Run("error", func(t *testing.T) {
		t.Run("prepare", func(t *testing.T) {
			layers := []Exec{
				&execObj{},
				&execObj{},
				&execObj{
					mockFailPrepare: true, mockFailCheck: true,
					mockFailRun: true, mockFailClose: true,
				},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := [][]int{
				[]int{1, 1, 0, 1},
				[]int{1, 1, 0, 1},
				[]int{1, 0, 0, 1},
			}

			// common errors
			switch err := run(layers...); {
			case err == nil:
				t.Fatal("Expected error, got nil")
			case err.Error() != "prepare failed":
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			for i, obj := range layers {
				if obj.(*execObj).countPrepare != flags[i][0] {
					t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
				}
				if obj.(*execObj).countCheck != flags[i][1] {
					t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
				}
				if obj.(*execObj).countRun != flags[i][2] {
					t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
				}
				if obj.(*execObj).countClose != flags[i][3] {
					t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
				}
			}
		})

		t.Run("check", func(t *testing.T) {
			layers := []Exec{
				&execObj{mockFailCheck: true},
				&execObj{},
				&execObj{
					mockFailPrepare: true, mockFailCheck: true,
					mockFailRun: true, mockFailClose: true,
				},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := [][]int{
				[]int{1, 1, 0, 1},
				[]int{0, 0, 0, 0},
				[]int{0, 0, 0, 0},
			}

			// common errors
			switch err := run(layers...); {
			case err == nil:
				t.Fatal("Expected error, got nil")
			case err.Error() != "check failed":
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			for i, obj := range layers {
				if obj.(*execObj).countPrepare != flags[i][0] {
					t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
				}
				if obj.(*execObj).countCheck != flags[i][1] {
					t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
				}
				if obj.(*execObj).countRun != flags[i][2] {
					t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
				}
				if obj.(*execObj).countClose != flags[i][3] {
					t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
				}
			}
		})

		t.Run("run", func(t *testing.T) {
			layers := []Exec{
				&execObj{},
				&execObj{mockFailRun: true},
				&execObj{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := [][]int{
				[]int{1, 1, 1, 1},
				[]int{1, 1, 1, 1},
				[]int{1, 1, 1, 1},
			}

			// common errors
			switch err := run(layers...); {
			case err == nil:
				t.Fatal("Expected error, got nil")
			case err.Error() != "run failed":
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			for i, obj := range layers {
				if obj.(*execObj).countPrepare != flags[i][0] {
					t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
				}
				if obj.(*execObj).countCheck != flags[i][1] {
					t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
				}
				if obj.(*execObj).countRun != flags[i][2] {
					t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
				}
				if obj.(*execObj).countClose != flags[i][3] {
					t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
				}
			}
		})

		t.Run("close", func(t *testing.T) {
			layers := []Exec{
				&execObj{},
				&execObj{mockFailClose: true},
				&execObj{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			flags := [][]int{
				[]int{1, 1, 1, 1},
				[]int{1, 1, 1, 1},
				[]int{1, 1, 1, 1},
			}

			// common errors
			switch err := run(layers...); {
			case err == nil:
				t.Fatal("Expected error, got nil")
			case err.Error() != "close failed":
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			for i, obj := range layers {
				if obj.(*execObj).countPrepare != flags[i][0] {
					t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
				}
				if obj.(*execObj).countCheck != flags[i][1] {
					t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
				}
				if obj.(*execObj).countRun != flags[i][2] {
					t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
				}
				if obj.(*execObj).countClose != flags[i][3] {
					t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
				}
			}
		})
	})

	t.Run("pipeline", func(t *testing.T) {
		t.Run("fake", func(t *testing.T) {
			layers := []Layer{
				&execObj{},
			}

			pipeline := &Pipeline{
				layers: layers,
			}

			flags := [][]int{
				[]int{1, 1, 1, 1},
			}

			// common errors
			run(pipeline)
			// switch err := run(layers...); {
			// case err == nil:
			// 	t.Fatal("Expected error, got nil")
			// case err.Error() != "close failed":
			// 	t.Fatalf("Unexpected error, got %q", err.Error())
			// }

			// table tests
			for i, obj := range layers {
				if obj.(*execObj).countPrepare != flags[i][0] {
					t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*execObj).countPrepare)
				}
				if obj.(*execObj).countCheck != flags[i][1] {
					t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*execObj).countCheck)
				}
				if obj.(*execObj).countRun != flags[i][2] {
					t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*execObj).countRun)
				}
				if obj.(*execObj).countClose != flags[i][3] {
					t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*execObj).countClose)
				}
			}
		})

		t.Run("real", func(t *testing.T) {
			// process1 := NewProcess("cat", "/dev/stdin")    // read 'foobar' from stdin
			// 	process2 := NewProcess("rev")                  // reverse -> raboof
			// 	process3 := NewProcess("grep", "-o", "raboof") // grep reversed (must be 1 match)
			// 	process4 := NewProcess("wc", "-l")             // count matches
		})
	})
}
