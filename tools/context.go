package tools

import (
	"bytes"
	"context"
	"text/template"

	srvctx "github.com/boomfunc/base/server/context"
)

func render(ctx context.Context, s string) string {
	q := func(key string) (string, error) { return srvctx.GetQ(ctx, key) }
	url := func(key string) (string, error) { return srvctx.GetUrl(ctx, key) }
	meta := func(key string) (interface{}, error) { return srvctx.GetMeta(ctx, key) }

	// add shortcuts to template rendering
	funcMap := template.FuncMap{
		"meta": meta,
		"q":    q,
		"url":  url,
	}
	tpl, err := template.New("").Funcs(funcMap).Parse(s)
	if err != nil {
		return s
	}

	// TODO look for better solution
	var b bytes.Buffer

	if err := tpl.Execute(&b, nil); err != nil {
		return s
	}

	return b.String()
}

func StringFromCtx(ctx context.Context, s string) string {
	return render(ctx, s)
}

// func FillStrings(ctx context.Context, ss []string) []string {
// 	for i, s := range ss {
// 		ns, _ := render(ctx, s)
// 		ss[i] = ns
// 	}
// 	return ss
// }
