package handler

import (
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/MerlinDMC/dsapid/storage"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"net/http"
)

func ApiPostDatasetUpdate(encoder middleware.OutputEncoder, params martini.Params, manifests storage.ManifestStorage, converter converter.DsapiManifestEncoder, user middleware.User, req *http.Request) (int, []byte) {
	action := req.URL.Query().Get("action")

	if manifest, ok := manifests.GetOK(params["id"]); ok && action != "" {
		switch action {
		case "enable":
			log.WithFields(log.Fields{
				"user_uuid":     user.GetId(),
				"user_name":     user.GetName(),
				"image_uuid":    manifest.Uuid,
				"image_name":    manifest.Name,
				"image_version": manifest.Version,
			}).Info("enabling image")

			manifest.State = dsapid.ManifestStateActive
			manifest.Disabled = false
			break
		case "deprecate":
			log.WithFields(log.Fields{
				"user_uuid":     user.GetId(),
				"user_name":     user.GetName(),
				"image_uuid":    manifest.Uuid,
				"image_name":    manifest.Name,
				"image_version": manifest.Version,
			}).Info("deprecating image")

			manifest.State = dsapid.ManifestStateDeprecated
			manifest.Disabled = false
			break
		case "disable":
			log.WithFields(log.Fields{
				"user_uuid":     user.GetId(),
				"user_name":     user.GetName(),
				"image_uuid":    manifest.Uuid,
				"image_name":    manifest.Name,
				"image_version": manifest.Version,
			}).Info("disabling image")

			manifest.State = dsapid.ManifestStateDisabled
			manifest.Disabled = true
			break
		case "nuke":
			log.WithFields(log.Fields{
				"user_uuid":     user.GetId(),
				"user_name":     user.GetName(),
				"image_uuid":    manifest.Uuid,
				"image_name":    manifest.Name,
				"image_version": manifest.Version,
			}).Info("nuking image")

			manifest.State = dsapid.ManifestStateNuked
			manifest.Disabled = true
			break
		}

		if err := manifests.Update(manifest.Uuid, manifest); err == nil {
			return http.StatusOK, encoder.MustEncode(converter.EncodeWithExtra(manifest))
		}
	}

	return http.StatusInternalServerError, encoder.MustEncode(dsapid.Table{
		"error": "update failed",
	})
}

func ApiPostReloadDatasets(encoder middleware.OutputEncoder, manifests storage.ManifestStorage, user middleware.User, req *http.Request) (int, []byte) {
	log.WithFields(log.Fields{
		"user_uuid": user.GetId(),
		"user_name": user.GetName(),
	}).Info("reloading datasets")

	manifests.Reload()

	return http.StatusOK, encoder.MustEncode(dsapid.Table{
		"ok": "datasets reloaded",
	})
}
