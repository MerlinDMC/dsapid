package imgapi

import (
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/storage"
	"net/url"
)

type imgapiEncoder struct {
	base  *url.URL
	users storage.UserStorage
}

func NewEncoder(base_url string, users storage.UserStorage) (encoder *imgapiEncoder) {
	encoder = new(imgapiEncoder)

	if base, err := url.Parse(base_url); err == nil {
		encoder.base = base
	}

	encoder.users = users

	return encoder
}

func (me *imgapiEncoder) Encode(manifest *dsapid.ManifestResource) interface{} {
	out := imgapiManifest{
		V:            CurrentManifestVersion,
		Uuid:         manifest.Uuid,
		Name:         manifest.Name,
		Version:      manifest.Version,
		Description:  manifest.Description,
		Os:           manifest.Os,
		Type:         manifest.Type,
		Urn:          manifest.Urn,
		State:        manifest.State,
		Disabled:     manifest.Disabled,
		Public:       manifest.Public,
		PublishedAt:  manifest.PublishedAt,
		Requirements: manifest.Requirements,
		Users:        manifest.Users,
		Tags:         manifest.Tags,
	}

	if owner, ok := me.users.GetOK(manifest.Owner); ok {
		out.Owner = owner.Uuid
	}

	if manifest.Type == dsapid.ManifestTypeZvol {
		if v, ok := manifest.Options["cpu_type"]; ok {
			out.CpuType = v.(string)
		}

		if v, ok := manifest.Options["image_size"]; ok {
			out.ImageSize = converter.DecodeToInt64(v)
		}

		if v, ok := manifest.Options["nic_driver"]; ok {
			out.NicDriver = v.(string)
		}

		if v, ok := manifest.Options["disk_driver"]; ok {
			out.DiskDriver = v.(string)
		}
	}

	for _, file := range manifest.Files {
		out.Files = append(out.Files, imgapiManifestFile{
			Size:        file.Size,
			Md5:         file.Md5,
			Sha1:        file.Sha1,
			Compression: file.Compression,
		})
	}

	return out
}

func (me *imgapiEncoder) EncodeWithExtra(manifest *dsapid.ManifestResource) interface{} {
	return me.Encode(manifest)
}
