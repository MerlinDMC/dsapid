package decoder

import (
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/converter/dsapi"
	"github.com/MerlinDMC/dsapid/converter/imgapi"
	"github.com/MerlinDMC/dsapid/storage"
)

func DecodeToManifest(data dsapid.Table, provider dsapid.SyncProvider, users storage.UserStorage) *dsapid.ManifestResource {
	var decoder converter.ManifestDecoder

	if _, ok := data["v"]; ok {
		decoder, _ = imgapi.NewDecoder(provider, users)
	} else {
		decoder, _ = dsapi.NewDecoder(provider, users)
	}

	return decoder.Decode(data)
}
