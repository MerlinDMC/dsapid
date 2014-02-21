package handler

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/MerlinDMC/dsapid/storage"
	"github.com/codegangsta/martini"
	"net/http"
)

func ImgapiList(encoder middleware.OutputEncoder, params martini.Params, manifests storage.ManifestStorage, converter converter.ImgapiManifestEncoder, user middleware.User, req *http.Request) (int, []byte) {
	var data []interface{} = make([]interface{}, 0)
	var filters []storage.ManifestFilter = []storage.ManifestFilter{storage.FilterManifestEnabled()}

	if user.IsGuest() {
		filters = append(filters, storage.FilterManifestPublic(true))
	} else {
		filters = append(filters, storage.FilterManifestForUser(user.GetId()))
	}

	if v := req.URL.Query().Get("name"); v != "" {
		filters = append(filters, storage.FilterManifestName(v))
	}

	if v := req.URL.Query().Get("version"); v != "" {
		filters = append(filters, storage.FilterManifestVersion(v))
	}

	if v := req.URL.Query().Get("os"); v != "" {
		filters = append(filters, storage.FilterManifestOs(v))
	}

	for manifest := range manifests.Filter(filters...) {
		data = append(data, converter.Encode(manifest))
	}

	return http.StatusOK, encoder.MustEncode(data)
}

func ImgapiDetail(encoder middleware.OutputEncoder, params martini.Params, manifests storage.ManifestStorage, converter converter.ImgapiManifestEncoder) (int, []byte) {
	if manifest, ok := manifests.GetOK(params["id"]); ok {
		return http.StatusOK, encoder.MustEncode(converter.Encode(manifest))
	}

	return http.StatusNotFound, []byte("Not found")
}

func ImgapiFile(params martini.Params, manifests storage.ManifestStorage, res http.ResponseWriter, req *http.Request) {
	if manifest, ok := manifests.GetOK(params["id"]); ok {
		var file_idx int = 0

		if v, ok := params["file_idx"]; ok {
			file_idx = int(converter.DecodeToInt64(v))
		}

		if len(manifest.Files) > file_idx {
			file := manifest.Files[file_idx]
			if md5_sum, err := hex.DecodeString(file.Md5); err == nil {
				res.Header().Set("Content-Type", "application/octet-stream")
				res.Header().Set("Content-Md5", base64.StdEncoding.EncodeToString(md5_sum))

				http.ServeFile(res, req, manifests.FilePath(manifest, &file))

				return
			}
		}
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte("Not found"))
}
