package executor

import (
	"context"
	"testing"
)

func TestOperation(t *testing.T) {
	obj := &fake{}

	op := Operation(
		[]OperationFunc{obj.a, obj.b},
		[]OperationFunc{obj.c},
	)

	if len(op.up) != 2 {
		t.Fatalf("len(op.up): expected \"%d\", got \"%d\"", 2, len(op.up))
	}
	if len(op.down) != 1 {
		t.Fatalf("len(op.down): expected \"%d\", got \"%d\"", 1, len(op.down))
	}
}

func TestOperationRun(t *testing.T) {
	ctx := context.TODO()

	t.Run("success", func(t *testing.T) {
		t.Run("up-/down-", func(t *testing.T) {
			objs := []*fake{
				&fake{},
			}
			matrix := [][]int{
				[]int{0, 0, 0}, // nothing to forward, nothing to backwards
			}
			// run action
			op := Operation(
				nil,
				nil,
			)
			if err := op.Run(ctx); err != nil {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}
			// check for matrix state
			checkMatrix(t, objs, matrix)
		})

		t.Run("up-/down+", func(t *testing.T) {
			objs := []*fake{
				&fake{},
			}
			matrix := [][]int{
				[]int{0, 0, 0}, // nothing to forward, nothing to backwards
			}
			// run action
			op := Operation(
				nil,
				[]OperationFunc{objs[0].a, objs[0].b},
			)
			if err := op.Run(ctx); err != nil {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}
			// check for matrix state
			checkMatrix(t, objs, matrix)
		})

		t.Run("up+/down-", func(t *testing.T) {
			objs := []*fake{
				&fake{},
			}
			matrix := [][]int{
				[]int{1, 1, 0}, // 2 operations to forward, nothing to backwards
			}
			// run action
			op := Operation(
				[]OperationFunc{objs[0].a, objs[0].b},
				nil,
			)
			if err := op.Run(ctx); err != nil {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}
			// check for matrix state
			checkMatrix(t, objs, matrix)
		})

		t.Run("up+/down+", func(t *testing.T) {
			objs := []*fake{
				&fake{},
			}
			matrix := [][]int{
				[]int{1, 1, 0}, // TODO maybe wrong logic? 2 operations to forward, but no error -> no operations for rollback
			}
			// run action
			op := Operation(
				[]OperationFunc{objs[0].a, objs[0].b},
				[]OperationFunc{objs[0].c, objs[0].b},
			)
			if err := op.Run(ctx); err != nil {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			// check for matrix state
			checkMatrix(t, objs, matrix)
		})
	})

	t.Run("error", func(t *testing.T) {
		t.Run("up+/down-", func(t *testing.T) {
			objs := []*fake{
				&fake{mockAA: true},
			}
			matrix := [][]int{
				[]int{1, 0, 0}, // nothing to forward, nothing to backwards
			}
			// run action
			op := Operation(
				[]OperationFunc{objs[0].a},
				nil,
			)
			err := op.Run(ctx)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != "a failed" {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}
			// check for matrix state
			checkMatrix(t, objs, matrix)
		})

		t.Run("up+/down+", func(t *testing.T) {
			objs := []*fake{
				&fake{mockAA: true},
			}
			matrix := [][]int{
				[]int{1, 1, 1}, // nothing to forward, nothing to backwards
			}
			// run action
			op := Operation(
				[]OperationFunc{objs[0].a},
				[]OperationFunc{objs[0].b, objs[0].c},
			)
			err := op.Run(ctx)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != "a failed" {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}
			// check for matrix state
			checkMatrix(t, objs, matrix)
		})
	})
}
