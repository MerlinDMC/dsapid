package middleware

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
	"strings"
)

type OutputEncoder interface {
	Encode(v ...interface{}) ([]byte, error)
	MustEncode(v ...interface{}) []byte
}

type jsonEncoder struct{}

func (me jsonEncoder) Encode(v ...interface{}) ([]byte, error) {
	var data interface{} = v

	if v == nil {
		// So that empty results produces `[]` and not `null`
		data = []interface{}{}
	} else if len(v) == 1 {
		data = v[0]
	}

	return json.Marshal(data)
}

func (me jsonEncoder) MustEncode(v ...interface{}) []byte {
	if data, err := me.Encode(v...); err == nil {
		return data
	}

	return nil
}

type prettyJsonEncoder struct{}

func (me prettyJsonEncoder) Encode(v ...interface{}) ([]byte, error) {
	var data interface{} = v

	if v == nil {
		// So that empty results produces `[]` and not `null`
		data = []interface{}{}
	} else if len(v) == 1 {
		data = v[0]
	}

	return json.MarshalIndent(data, "", "  ")
}

func (me prettyJsonEncoder) MustEncode(v ...interface{}) []byte {
	if data, err := me.Encode(v...); err == nil {
		return data
	}

	return nil
}

func EncodeOutput(prettify bool) martini.Handler {
	return func(ctx martini.Context, res http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.RequestURI, "/ui") {
			return
		}

		if prettify {
			ctx.MapTo(prettyJsonEncoder{}, (*OutputEncoder)(nil))
		} else {
			ctx.MapTo(jsonEncoder{}, (*OutputEncoder)(nil))
		}

		res.Header().Set("Content-Type", "application/json")
	}
}
