package executor

import (
	"context"
	"testing"
)

func TestExecutor(t *testing.T) {
	obj := &fake{mockBB: true}

	ex := New(
		context.TODO(),
		Operation(
			[]OperationFunc{obj.a},
			nil,
			false,
		),
	)

	t.Run("New", func(t *testing.T) {
		if n := len(ex.operations); n != 1 {
			t.Fatalf("len(ex.operations): expected \"%d\", got \"%d\"", 1, n)
		}
	})

	t.Run("AddOperations", func(t *testing.T) {
		ex.AddOperations(
			Operation(
				[]OperationFunc{func(ctx context.Context) error { return nil }},
				[]OperationFunc{func(ctx context.Context) error { return nil }},
				false,
			),
		)

		if n := len(ex.operations); n != 2 {
			t.Fatalf("len(ex.operations): expected \"%d\", got \"%d\"", 2, n)
		}
	})

	t.Run("Run", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			if err := ex.Run(); err != nil {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			matrix := [][]int{
				[]int{1, 0, 0},
			}
			checkMatrix(t, []*fake{obj}, matrix)
		})

		t.Run("Error", func(t *testing.T) {
			ex.AddOperations(
				Operation(
					[]OperationFunc{obj.b},
					[]OperationFunc{obj.c},
					false,
				),
			)

			err := ex.Run()
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != "b failed" {
				t.Fatalf("Unexpected error, got %q", err.Error())
			}

			matrix := [][]int{
				[]int{2, 1, 1}, // a called twice because in Run/Success was first call
			}
			checkMatrix(t, []*fake{obj}, matrix)
		})
	})
}
