package executor

import (
	"testing"
	"context"
)

type fake struct {
	aa, bb, cc int
}

func (f *fake) a(ctx context.Context) error {
	f.aa++
	return nil
}
func (f *fake) b(ctx context.Context) error {
	f.bb++
	return nil
}
func (f *fake) c(ctx context.Context) error {
	f.cc++
	return nil
}

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

func TestConcurrent(t *testing.T) {
	t.Run("initial", func(t *testing.T) {
		objs := []*fake{
			&fake{0, 0, 0},
			&fake{0, 1, 0},
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
}
