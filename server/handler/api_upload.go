package handler

import (
	"encoding/json"
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter/decoder"
	"github.com/MerlinDMC/dsapid/server/logger"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/MerlinDMC/dsapid/storage"
	"github.com/codegangsta/martini"
	"io"
	"net/http"
	"os"
	"time"
)

func ApiPostFileUpload(encoder middleware.OutputEncoder, params martini.Params, manifests storage.ManifestStorage, users storage.UserStorage, user middleware.User, req *http.Request) (int, []byte) {
	var manifest *dsapid.ManifestResource

	if file, _, err := req.FormFile("manifest"); err == nil {
		var data dsapid.Table

		if err := json.NewDecoder(file).Decode(&data); err == nil {
			manifest = decoder.DecodeToManifest(data, dsapid.SyncProviderCommunity, users)

			if _, ok := manifests.GetOK(manifest.Uuid); ok {
				logger.Infof("uploading duplicate image: %s (%s v%s)", manifest.Uuid, manifest.Name, manifest.Version)

				return http.StatusInternalServerError, encoder.MustEncode(dsapid.Table{
					"error": "image already exists",
				})
			}

			logger.Infof("uploading image: %s (%s v%s)", manifest.Uuid, manifest.Name, manifest.Version)

			manifest.PublishedAt = time.Now()
			manifest.State = dsapid.ManifestStatePending

			if !user.HasRoles(dsapid.UserRoleDatasetManage) && !user.HasRoles(dsapid.UserRoleDatasetAdmin) {
				manifest.Owner = user.GetId()
			}

			if user.HasRoles(dsapid.UserRoleDatasetAdmin) {
				manifest.State = dsapid.ManifestStateActive
			}

			if file, _, err := req.FormFile("file"); err == nil {
				if err = os.MkdirAll(manifests.ManifestPath(manifest), 0770); err == nil {
					if file_out, err := os.Create(manifests.FilePath(manifest, &manifest.Files[0])); err == nil {
						if n, err := io.Copy(file_out, file); err == nil && n == manifest.Files[0].Size {
							manifests.Add(manifest.Uuid, manifest)

							return http.StatusOK, encoder.MustEncode(manifest)
						}
					}
				}
			}
		}
	}

	return http.StatusInternalServerError, encoder.MustEncode(dsapid.Table{
		"error": "upload failed",
	})
}
