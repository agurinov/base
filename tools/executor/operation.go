package executor

import (
	"context"
)

// operation is
// forward functions with rollback functions
type operation struct {
	up            []OperationFunc
	down          []OperationFunc
	forceRollback bool
}

func (op *operation) backward(ctx context.Context) error {
	if op.down != nil {
		return concurrent(ctx, op.down...)
	}
	return nil
}

func Operation(up, down []OperationFunc, forceRollback bool) *operation {
	return &operation{up, down, forceRollback}
}

// Run runs []OperationFunc associated with operation
// if up returns error - down is triggered
func (op *operation) Run(ctx context.Context) error {
	if op.up != nil {
		// forward stage
		if err := concurrent(ctx, op.up...); err != nil {
			// forward movement failed, need to rollback
			// but return error from run
			op.backward(ctx)
			return err
		} else if op.forceRollback {
			// forward movement failed, need to rollback
			// return error from bacward stage if exists
			return op.backward(ctx)
		}
	}

	return nil
}
