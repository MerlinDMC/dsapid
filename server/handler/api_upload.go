package handler

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter/decoder"
	"github.com/MerlinDMC/dsapid/server/logger"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/MerlinDMC/dsapid/storage"
	"github.com/go-martini/martini"
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
						defer file_out.Close()

						hash_md5 := md5.New()
						hash_sha1 := sha1.New()

						writer := io.MultiWriter(hash_md5, hash_sha1, file_out)

						if _, err := io.Copy(writer, file); err == nil {
							md5_sum := hex.EncodeToString(hash_md5.Sum(nil))
							sha1_sum := hex.EncodeToString(hash_sha1.Sum(nil))

							if manifest.Files[0].Md5 != "" && manifest.Files[0].Md5 != md5_sum {
								logger.Warnf("checksum mismatch: got %s expected %s", md5_sum, manifest.Files[0].Md5)
								goto errCancel
							}

							if manifest.Files[0].Sha1 != "" && manifest.Files[0].Sha1 != sha1_sum {
								logger.Warnf("checksum mismatch: got %s expected %s", sha1_sum, manifest.Files[0].Sha1)
								goto errCancel
							}

							manifest.Files[0].Md5 = md5_sum
							manifest.Files[0].Sha1 = sha1_sum

							manifests.Add(manifest.Uuid, manifest)

							return http.StatusOK, encoder.MustEncode(manifest)
						}
					}
				}
			}
		}
	}

errCancel:
	if manifest.Uuid != "" {
		manifests.Delete(manifest.Uuid)
	}

	return http.StatusInternalServerError, encoder.MustEncode(dsapid.Table{
		"error": "upload failed",
	})
}
