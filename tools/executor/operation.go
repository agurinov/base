package executor

import (
	"context"
)

// operation is
// forward functions with rollback functions
type operation struct {
	up   []OperationFunc
	down []OperationFunc
}

func Operation(up, down []OperationFunc) *operation {
	return &operation{up, down}
}

// Run runs []OperationFunc associated with operation
// if up returns error - down is triggered
func (op *operation) Run(ctx context.Context) error {
	if op.up != nil {
		// forward stage
		if err := concurrent(ctx, op.up...); err != nil {
			// forward movement failed, need to rollback
			if op.down != nil {
				// rollback
				concurrent(ctx, op.down...)
			}
			return err
		}
	}

	return nil
}
