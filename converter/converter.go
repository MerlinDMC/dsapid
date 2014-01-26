package converter

import (
	"fmt"
	"github.com/MerlinDMC/dsapid"
	"strconv"
	"time"
)

type ManifestEncoder interface {
	Encode(*dsapid.ManifestResource) interface{}
	EncodeWithExtra(*dsapid.ManifestResource) interface{}
}

type DsapiManifestEncoder interface {
	ManifestEncoder
}

type ImgapiManifestEncoder interface {
	ManifestEncoder
}

type ManifestDecoder interface {
	Decode(dsapid.Table) *dsapid.ManifestResource
}

func DecodeToInt64(v interface{}) int64 {
	switch v.(type) {
	case int64:
		return v.(int64)
		break
	case int:
		return int64(v.(int))
		break
	case float64:
		return int64(v.(float64))
		break
	case string:
		if i, err := strconv.ParseInt(v.(string), 10, 0); err == nil {
			return i
		}
		break
	}

	return -1
}

func ComputeUrn(manifest *dsapid.ManifestResource) string {
	return fmt.Sprintf("smartos:smartos:%s:%s", manifest.Name, manifest.Version)
}

func ParseDateTime(value string) (converted time.Time, err error) {
	var formats = []string{
		"2006-01-02T15:04:05.999999999Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04Z",
		"2006-01-02T15:04:05", // old datasets.at
	}

	for _, fmt := range formats {
		if converted, err = time.Parse(fmt, value); err == nil {
			return
		}
	}

	return
}
