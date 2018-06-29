package conf

// import (
// 	"bytes"
// 	"testing"
//
// 	"app/pipeline"
// )
//
// func TestRouteRunParallel(t *testing.T) {
// 	layers := []pipeline.Layer{
// 		&pipeline.FakeLayer{},
// 	}
//
// 	pipeline := pipeline.New(layers)
//
// 	route := Route{
// 		pipeline: pipeline,
// 	}
//
// 	input := bytes.NewBuffer([]byte{})
// 	output := bytes.NewBuffer([]byte{})
//
// 	t.Run("first", func(t *testing.T) {
// 		t.Parallel()
//
// 		route.Run(input, output)
// 	})
//
// 	t.Run("second", func(t *testing.T) {
// 		t.Parallel()
//
// 		route.Run(input, output)
// 	})
// }
