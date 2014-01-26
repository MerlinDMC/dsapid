package dsapi

import (
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/storage"
	"net/url"
	"path"
)

type dsapiEncoder struct {
	base  *url.URL
	users storage.UserStorage
}

func NewEncoder(base_url string, users storage.UserStorage) (encoder *dsapiEncoder) {
	encoder = new(dsapiEncoder)

	if base, err := url.Parse(base_url); err == nil {
		encoder.base = base
	}

	encoder.users = users

	return encoder
}

func (me *dsapiEncoder) Encode(manifest *dsapid.ManifestResource) interface{} {
	out := dsapiManifest{
		Uuid:         manifest.Uuid,
		Name:         manifest.Name,
		Version:      manifest.Version,
		Description:  manifest.Description,
		Os:           manifest.Os,
		Type:         manifest.Type,
		Urn:          manifest.Urn,
		PublishedAt:  manifest.PublishedAt,
		CreatedAt:    manifest.CreatedAt,
		Requirements: manifest.Requirements,
		Users:        manifest.Users,
		Tags:         manifest.Tags,
	}

	if owner, ok := me.users.GetOK(manifest.Owner); ok {
		out.CreatorUuid = owner.Uuid
		out.CreatorName = owner.Name
		out.VendorUuid = owner.Uuid
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

	for file_idx, file := range manifest.Files {
		out.Files = append(out.Files, dsapiManifestFile{
			Path: file.Path,
			Size: file.Size,
			Md5:  file.Md5,
			Sha1: file.Sha1,
		})

		if file_url, err := url.Parse(path.Join("", "datasets", out.Uuid, file.Path)); err == nil {
			out.Files[file_idx].Url = me.base.ResolveReference(file_url).String()
		}
	}

	return out
}

func (me *dsapiEncoder) EncodeWithExtra(manifest *dsapid.ManifestResource) interface{} {
	out := me.Encode(manifest).(dsapiManifest)

	out.BuilderInfo = manifest.BuilderInfo
	out.MetadataInfo = manifest.MetadataInfo

	return out
}
