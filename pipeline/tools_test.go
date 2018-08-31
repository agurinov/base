package pipeline

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/boomfunc/base/tools"
)

// check object that not started for flow
// if prepare == 0
func checkNullMatrix(t *testing.T, obj Exec, i int) {
	o, ok := obj.(*FakeLayer)
	if !ok {
		t.Fatalf("objs[%d]: inappropriate type", i)
	}

	if o.countPrepare != 0 {
		fmt.Println("DMLDMLKD:", o)
		t.Fatalf("objs[%d].countPrepare: expected \"%d\", got \"%d\"", i, 0, o.countPrepare)
	}
	if o.countCheck != 0 {
		fmt.Println("DMLDMLKD:", o)
		t.Fatalf("objs[%d].countCheck: expected \"%d\", got \"%d\"", i, 0, o.countCheck)
	}
	if o.countRun != 0 {
		fmt.Println("DMLDMLKD:", o)
		t.Fatalf("objs[%d].countRun: expected \"%d\", got \"%d\"", i, 0, o.countRun)
	}
	if o.countClose != 0 {
		fmt.Println("DMLDMLKD:", o)
		t.Fatalf("objs[%d].countClose: expected \"%d\", got \"%d\"", i, 0, o.countClose)
	}
}

func checkMatrix(t *testing.T, matrix [][]int, objs []Exec) {
	// check sizes
	if len(matrix) != len(objs) {
		t.Fatalf("CheckByMatrix: inappropriate sizes, len(matrix)==%d != len(objs)==%d", len(matrix), len(objs))
	}

	// table tests
	for i, obj := range objs {
		o, ok := obj.(*FakeLayer)
		if !ok {
			t.Fatalf("objs[%d]: inappropriate type", i)
		}

		if o.countPrepare == 0 {
			checkNullMatrix(t, o, i)
			return
		}

		if o.countPrepare != matrix[i][0] {
			fmt.Println("DMLDMLKD:", o)
			t.Fatalf("objs[%d].countPrepare: expected \"%d\", got \"%d\"", i, matrix[i][0], o.countPrepare)
		}
		if o.countCheck != matrix[i][1] {
			fmt.Println("DMLDMLKD:", o)
			t.Fatalf("objs[%d].countCheck: expected \"%d\", got \"%d\"", i, matrix[i][1], o.countCheck)
		}
		if o.countRun != matrix[i][2] {
			fmt.Println("DMLDMLKD:", o)
			t.Fatalf("objs[%d].countRun: expected \"%d\", got \"%d\"", i, matrix[i][2], o.countRun)
		}
		if o.countClose != matrix[i][3] {
			fmt.Println("DMLDMLKD:", o)
			t.Fatalf("objs[%d].countClose: expected \"%d\", got \"%d\"", i, matrix[i][3], o.countClose)
		}
	}
}

func TestPiping(t *testing.T) {
	input := tools.ReadCloser(bytes.NewBuffer([]byte{}))
	output := tools.WriteCloser(bytes.NewBuffer([]byte{}))
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

func TestRun(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("len==1", func(t *testing.T) {
			layers := []Exec{
				&FakeLayer{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			matrix := [][]int{
				[]int{1, 1, 1, 1},
			}

			// common errors
			if err := run(context.TODO(), layers...); err != nil {
				t.Fatal(err)
			}

			// table tests
			checkMatrix(t, matrix, layers)
		})

		t.Run("len==3", func(t *testing.T) {
			layers := []Exec{
				&FakeLayer{},
				&FakeLayer{},
				&FakeLayer{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			matrix := [][]int{
				[]int{1, 1, 1, 1},
				[]int{1, 1, 1, 1},
				[]int{1, 1, 1, 1},
			}

			// common errors
			if err := run(context.TODO(), layers...); err != nil {
				t.Fatal(err)
			}

			// table tests
			checkMatrix(t, matrix, layers)
		})
	})

	t.Run("error", func(t *testing.T) {
		t.Run("prepare", func(t *testing.T) {
			layers := []Exec{
				&FakeLayer{},
				&FakeLayer{},
				&FakeLayer{mockFailPrepare: true},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			matrix := [][]int{
				[]int{0, 0, 0, 0},
				[]int{0, 0, 0, 0},
				[]int{1, 0, 0, 1},
			}

			// common errors
			switch err := run(context.TODO(), layers...); {
			case err == nil:
				t.Fatal("Expected error, got nil")
			case err.Error() != "prepare failed":
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			checkMatrix(t, matrix, layers)
		})

		t.Run("check", func(t *testing.T) {
			layers := []Exec{
				&FakeLayer{},
				&FakeLayer{mockFailCheck: true},
				&FakeLayer{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			matrix := [][]int{
				[]int{0, 0, 0, 0},
				[]int{1, 1, 0, 1},
				[]int{0, 0, 0, 0},
			}

			// common errors
			switch err := run(context.TODO(), layers...); {
			case err == nil:
				t.Fatal("Expected error, got nil")
			case err.Error() != "check failed":
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			checkMatrix(t, matrix, layers)
		})

		t.Run("run", func(t *testing.T) {
			layers := []Exec{
				&FakeLayer{},
				&FakeLayer{mockFailRun: true},
				&FakeLayer{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			matrix := [][]int{
				[]int{0, 0, 0, 0},
				[]int{1, 1, 1, 1},
				[]int{0, 0, 0, 0},
			}

			// common errors
			switch err := run(context.TODO(), layers...); {
			case err == nil:
				t.Fatal("Expected error, got nil")
			case err.Error() != "run failed":
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			checkMatrix(t, matrix, layers)
		})

		t.Run("close", func(t *testing.T) {
			layers := []Exec{
				&FakeLayer{},
				&FakeLayer{mockFailClose: true},
				&FakeLayer{},
			}
			// digits - count of invokes functions
			// digits from left to right:
			// countPrepare, countCheck, countRun, countClose
			matrix := [][]int{
				[]int{0, 0, 0, 0},
				[]int{1, 1, 1, 1},
				[]int{0, 0, 0, 0},
			}

			// common errors
			switch err := run(context.TODO(), layers...); {
			case err == nil:
				t.Fatal("Expected error, got nil")
			case err.Error() != "close failed":
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// table tests
			checkMatrix(t, matrix, layers)
		})
	})

	t.Run("pipeline", func(t *testing.T) {
		// t.Run("fake", func(t *testing.T) {
		// 	pipeline := New(&FakeLayer{})
		//
		// 	flags := [][]int{
		// 		[]int{1, 1, 1, 1},
		// 	}
		//
		// 	input := bytes.NewBuffer([]byte{})
		// 	output := bytes.NewBuffer([]byte{})
		//
		// 	pipeline.Run(input, output)
		//
		// 	// table tests
		// 	for i, obj := range *pipeline {
		// 		if obj.(*FakeLayer).countPrepare != flags[i][0] {
		// 			t.Fatalf("layers[%d].countPrepare: expected \"%d\", got \"%d\"", i, flags[i][0], obj.(*FakeLayer).countPrepare)
		// 		}
		// 		if obj.(*FakeLayer).countCheck != flags[i][1] {
		// 			t.Fatalf("layers[%d].countCheck: expected \"%d\", got \"%d\"", i, flags[i][1], obj.(*FakeLayer).countCheck)
		// 		}
		// 		if obj.(*FakeLayer).countRun != flags[i][2] {
		// 			t.Fatalf("layers[%d].countRun: expected \"%d\", got \"%d\"", i, flags[i][2], obj.(*FakeLayer).countRun)
		// 		}
		// 		if obj.(*FakeLayer).countClose != flags[i][3] {
		// 			t.Fatalf("layers[%d].countClose: expected \"%d\", got \"%d\"", i, flags[i][3], obj.(*FakeLayer).countClose)
		// 		}
		// 	}
		// })

		t.Run("real", func(t *testing.T) {
			// process1 := NewProcess("cat", "/dev/stdin")    // read 'foobar' from stdin
			// 	process2 := NewProcess("rev")                  // reverse -> raboof
			// 	process3 := NewProcess("grep", "-o", "raboof") // grep reversed (must be 1 match)
			// 	process4 := NewProcess("wc", "-l")             // count matches
		})
	})
}
