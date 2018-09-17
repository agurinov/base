package application

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/boomfunc/base/server/flow"
)

// Load test
// JS='{"url":"geo","input":"185.86.151.11"}'
// seq 1000 | xargs -n 1 -P 250 sh -c "echo '$JS' | nc playground.lo 8080"
type jsonPacker struct{}

func (packer *jsonPacker) Unpack(ctx context.Context, r io.Reader) (*flow.Request, error) {
	intermediate := struct {
		Url   string
		Input string
	}{}

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&intermediate); err != nil {
		return nil, err
	}

	return flow.NewRequest(
		intermediate.Url,
		strings.NewReader(intermediate.Input),
	)
}

func (packer *jsonPacker) Pack(r io.Reader, w io.Writer) (int64, error) {
	return io.Copy(w, r)
}
