package executor

import (
	"context"
)

type OperationFunc func(context.Context) error
