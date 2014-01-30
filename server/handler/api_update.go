package handler

import (
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/server/logger"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/MerlinDMC/dsapid/storage"
	"github.com/codegangsta/martini"
	"net/http"
)

func ApiPostDatasetUpdate(encoder middleware.OutputEncoder, params martini.Params, manifests storage.ManifestStorage, converter converter.DsapiManifestEncoder, user middleware.User, req *http.Request) (int, []byte) {
	action := req.URL.Query().Get("action")

	if manifest, ok := manifests.GetOK(params["id"]); ok && action != "" {
		switch action {
		case "enable":
			logger.Infof("enabling image %s (user=%s)", manifest.Uuid, user.GetId())

			manifest.State = dsapid.ManifestStateActive
			manifest.Disabled = false
			break
		case "disable":
			logger.Infof("disabling image %s (user=%s)", manifest.Uuid, user.GetId())

			manifest.State = dsapid.ManifestStateDisabled
			manifest.Disabled = true
			break
		case "nuke":
			logger.Infof("nuking image %s (user=%s)", manifest.Uuid, user.GetId())

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
