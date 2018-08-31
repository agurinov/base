package pipeline

import (
	"context"
	"errors"
	"io"
)

var (
	ErrUnknownLayer = errors.New("pipeline: Unknown layer type")
)

// TODO rename LayerCopier
type Cloneable interface {
	copy() Layer
}

// Able interface describes an object that can be associated with other objects by stdio
type Able interface {
	setStdin(reader io.ReadCloser)
	setStdout(writer io.WriteCloser)
}

// Exec interface describes objects that can be self checked and can be executable by Pipeline
// Common lifecycle is:
// 		1. prepare (initial preparations)
// 		2. check (check that object can be executed)
// 		3. Run
// 		4. close (clear all object's sensitive data for reuse this object)
type Exec interface {
	prepare(ctx context.Context) error
	check(ctx context.Context) error
	run(ctx context.Context) error
	close(ctx context.Context) error // TODO look at priority order errors
}

// Layer interface describes complex type of object that can be a part of Pipeline
// in the Pipeline has the following lifecycle:
// 		1. prepare (create internal variables (cmd, connection, etc) )
// 		2. check (check that layer piped and can be executed)
// 		3. run
// 		4. Close (close stdio and clear all layer's sensitive data for reuse)
type Layer interface {
	Able
	Exec
}
