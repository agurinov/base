package tools

// import (
// 	"github.com/boomfunc/base/server/context"
// 	"net/url"
// 	"testing"
// 	// "os/exec"
// )
//
// func TestRender(t *testing.T) {
// 	s := "-lang=ru -ip={{q \"ip\"}}"
// 	ctx := context.New()
//
// 	q := url.Values{}
// 	q.Set("ip", "someip")
//
// 	vs, _ := context.Values(ctx)
// 	vs.Q = q
//
// 	t.Log(render(ctx, s))
// }
//
// func TestFillStrings(t *testing.T) {
// 	ss := []string{"-lang=ru -ip={{q \"ip\"}}"}
// 	ctx := context.New()
//
// 	q := url.Values{}
// 	q.Set("ip", "someip")
//
// 	vs, _ := context.Values(ctx)
// 	vs.Q = q
//
// 	t.Log(FillStrings(ctx, ss))
// }

// func TestMklremfle(t *testing.T) {
// 	out, err := exec.Command("/go/src/github.com/boomfunc/base/geo/geoip", "-lang=ru", "-ip=185.86.151.11").Output()
// 	// out, err := exec.Command("ls", "-lah", "/go/src/github.com/boomfunc/base/geo").Output()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(string(out))
// }
