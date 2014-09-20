package handler

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/MerlinDMC/dsapid/storage"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"io"
	"net/http"
	"os"
)

func ApiDatasetsList(encoder middleware.OutputEncoder, params martini.Params, manifests storage.ManifestStorage, user middleware.User, req *http.Request) (int, []byte) {
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
		data = append(data, manifest)
	}

	return http.StatusOK, encoder.MustEncode(data)
}

func ApiDatasetsDetail(encoder middleware.OutputEncoder, params martini.Params, manifests storage.ManifestStorage, converter converter.DsapiManifestEncoder, user middleware.User, req *http.Request) (int, []byte) {
	if manifest, ok := manifests.GetOK(params["id"]); ok {
		return http.StatusOK, encoder.MustEncode(manifest)
	}

	return http.StatusNotFound, []byte("Not found")
}

func ApiDatasetExport(encoder middleware.OutputEncoder, params martini.Params, manifests storage.ManifestStorage, dsapi_converter converter.DsapiManifestEncoder, imgapi_converter converter.ImgapiManifestEncoder, user middleware.User, res http.ResponseWriter) {
	if manifest, ok := manifests.GetOK(params["id"]); ok {
		res.Header().Set("Content-Type", "application/octet-stream")
		res.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s-%s.tar\"", manifest.Name, manifest.Version))

		tw := tar.NewWriter(res)
		defer tw.Close()

		if manifest_buf, err := json.MarshalIndent(dsapi_converter.Encode(manifest), "", "  "); err == nil {
			tw.WriteHeader(&tar.Header{
				Name: fmt.Sprintf("%s-%s-dsapi.dsmanifest", manifest.Name, manifest.Version),
				Mode: 0666,
				Size: int64(len(manifest_buf)),
			})
			tw.Write(manifest_buf)
		}

		if manifest_buf, err := json.MarshalIndent(imgapi_converter.Encode(manifest), "", "  "); err == nil {
			tw.WriteHeader(&tar.Header{
				Name: fmt.Sprintf("%s-%s-imgapi.dsmanifest", manifest.Name, manifest.Version),
				Mode: 0666,
				Size: int64(len(manifest_buf)),
			})
			tw.Write(manifest_buf)
		}

		for _, file := range manifest.Files {
			filename := manifests.FilePath(manifest, &file)

			if fin, err := os.OpenFile(filename, os.O_RDONLY, 0660); err == nil {
				tw.WriteHeader(&tar.Header{
					Name: manifest.Files[0].Path,
					Mode: 0666,
					Size: manifest.Files[0].Size,
				})

				if _, err := io.Copy(tw, fin); err != nil {
					log.Errorf("failed to create tar streaming archive: %s", err)
				}
			}
		}

		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte("Not found"))
}
