package handler

import (
	"fmt"
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/MerlinDMC/dsapid/storage"
	"net/http"
	"runtime"
)

func CommonPing(encoder middleware.OutputEncoder, user middleware.User) (int, []byte) {
	pingResponse := map[string]interface{}{
		"ping":    "pong",
		"app":     dsapid.AppName,
		"version": dsapid.AppVersion,
		"dsapi":   true,
		"imgapi":  true,
	}

	if !user.IsGuest() {
		pingResponse["auth"] = user.GetAuthInfo()
	}

	return http.StatusOK, encoder.MustEncode(pingResponse)
}

func CommonStatus(encoder middleware.OutputEncoder, manifests storage.ManifestStorage) (int, []byte) {
	manifests_count, manifests_size := int64(0), int64(0)

	for manifest := range manifests.List() {
		manifests_count++

		for _, file := range manifest.Files {
			manifests_size += file.Size
		}
	}

	statusResponse := map[string]interface{}{
		"app":            fmt.Sprintf("%s/%s", dsapid.AppName, dsapid.AppVersion),
		"manifest_count": manifests_count,
		"manifest_size":  manifests_size,
		"num_goroutines": runtime.NumGoroutine(),
	}

	return http.StatusOK, encoder.MustEncode(statusResponse)
}
