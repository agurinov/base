package executor

import (
	"context"
	"errors"
	"testing"
)

type fake struct {
	aa, bb, cc             int
	mockAA, mockBB, mockCC bool
}

func (f *fake) a(ctx context.Context) error {
	f.aa++

	if f.mockAA {
		return errors.New("a failed")
	}

	return nil
}
func (f *fake) b(ctx context.Context) error {
	f.bb++

	if f.mockBB {
		return errors.New("b failed")
	}

	return nil
}
func (f *fake) c(ctx context.Context) error {
	f.cc++

	if f.mockCC {
		return errors.New("c failed")
	}

	return nil
}

// special tool for checking objects state by matrix
func checkMatrix(t *testing.T, objs []*fake, matrix [][]int) {
	// check sizes
	if len(matrix) != len(objs) {
		t.Fatalf("CheckMatrix: inappropriate sizes, len(matrix)==%d != len(objs)==%d", len(matrix), len(objs))
	}

	// table tests
	for i, obj := range objs {
		if obj.aa != matrix[i][0] {
			t.Fatalf("objs[%d].aa: expected \"%d\", got \"%d\"", i, matrix[i][0], obj.aa)
		}
		if obj.bb != matrix[i][1] {
			t.Fatalf("objs[%d].bb: expected \"%d\", got \"%d\"", i, matrix[i][1], obj.bb)
		}
		if obj.cc != matrix[i][2] {
			t.Fatalf("objs[%d].cc: expected \"%d\", got \"%d\"", i, matrix[i][2], obj.cc)
		}
	}
}

func TestExecute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			errCh := make(chan error, 1)
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			objs := []*fake{
				&fake{},
			}
			matrix := [][]int{
				[]int{1, 0, 0},
			}

			execute(objs[0].a, errCh, ctx, cancel)

			// check for matrix state
			checkMatrix(t, objs, matrix)

			// check for errors
			if len(errCh) != 0 {
				t.Fatalf("errCh: expected \"%d\", got \"%d\"", 0, len(errCh))
			}
		})

		t.Run("cancelled", func(t *testing.T) {
			errCh := make(chan error, 1)
			ctx, cancel := context.WithCancel(context.TODO())

			objs := []*fake{
				&fake{},
			}
			matrix := [][]int{
				[]int{0, 0, 0},
			}

			// imitate cancelling from another goroutine
			cancel()

			execute(objs[0].a, errCh, ctx, cancel)

			// check for matrix state
			checkMatrix(t, objs, matrix)

			// check for errors
			if len(errCh) != 0 {
				t.Fatalf("errCh: expected \"%d\", got \"%d\"", 0, len(errCh))
			}
		})
	})

	t.Run("error", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			errCh := make(chan error, 1)
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			objs := []*fake{
				&fake{mockAA: true},
			}
			matrix := [][]int{
				[]int{1, 0, 0},
			}

			execute(objs[0].a, errCh, ctx, cancel)

			// check for matrix state
			checkMatrix(t, objs, matrix)

			// check for errors
			if len(errCh) != 1 {
				t.Fatalf("errCh: expected \"%d\", got \"%d\"", 1, len(errCh))
			}
			// get error
			if err := <-errCh; err.Error() != "a failed" {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}
		})
	})
}

func TestConcurrent(t *testing.T) {
	t.Run("initial", func(t *testing.T) {
		objs := []*fake{
			&fake{aa: 0, bb: 0, cc: 0},
			&fake{aa: 0, bb: 1, cc: 0},
		}
		matrix := [][]int{
			[]int{0, 0, 0},
			[]int{0, 1, 0},
		}
		checkMatrix(t, objs, matrix)
	})

	t.Run("success", func(t *testing.T) {
		objs := []*fake{
			&fake{},
		}
		matrix := [][]int{
			[]int{1, 1, 1},
		}

		if err := concurrent(context.TODO(), objs[0].a, objs[0].b, objs[0].c); err != nil {
			t.Fatal(err)
		}

		checkMatrix(t, objs, matrix)
	})

	t.Run("error", func(t *testing.T) {
		// b method will return error
		objs := []*fake{
			&fake{mockBB: true},
		}
		// matrix := [][]int{
		// 	[]int{1, 1, 0},
		// }

		err := concurrent(context.TODO(), objs[0].a, objs[0].b, objs[0].c)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "b failed" {
			t.Fatalf("Unexpected error, got %q", err.Error())
		}

		// time to time different because concurrent
		// checkMatrix(t, objs, matrix)
	})
}
