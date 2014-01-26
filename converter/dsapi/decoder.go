package dsapi

import (
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/storage"
	"path"
)

type dsapiDecoder struct {
	provider dsapid.SyncProvider
	users    storage.UserStorage
}

func NewDecoder(provider dsapid.SyncProvider, users storage.UserStorage) (decoder *dsapiDecoder, err error) {
	decoder = new(dsapiDecoder)

	if provider != "" {
		decoder.provider = provider
	} else {
		decoder.provider = dsapid.SyncProviderCommunity
	}

	decoder.users = users

	return decoder, nil
}

func (me *dsapiDecoder) Decode(data dsapid.Table) *dsapid.ManifestResource {
	var manifest *dsapid.ManifestResource

	manifest = new(dsapid.ManifestResource)

	manifest.Options = make(dsapid.Table)
	manifest.SyncInfo = make(dsapid.Table)
	manifest.BuilderInfo = make(dsapid.Table)
	manifest.MetadataInfo = make([]dsapid.Table, 0)

	var creator_uuid string = dsapid.DefaultUserUuid
	var creator_name string = dsapid.DefaultUserName

	if v, ok := data["creator_uuid"]; ok && v.(string) != "" {
		creator_uuid = v.(string)
		// } else {
		// 	logger.Warn("manifest has no creator_uuid")
	}

	if v, ok := data["creator_name"]; ok && v.(string) != "" {
		creator_name = v.(string)
		// } else {
		// 	logger.Warn("manifest has no creator_name")
	}

	user := me.users.EnsureExists(creator_uuid, creator_name)

	manifest.Owner = user.Uuid

	if user.Provider != "" {
		manifest.Provider = user.Provider
	} else {
		manifest.Provider = me.provider
	}

	if v, ok := data["uuid"]; ok {
		manifest.Uuid = v.(string)
	}

	if v, ok := data["name"]; ok {
		manifest.Name = v.(string)
	}

	if v, ok := data["version"]; ok {
		manifest.Version = v.(string)
	}

	if v, ok := data["description"]; ok {
		manifest.Description = v.(string)
	}

	if v, ok := data["type"]; ok {
		manifest.Type = dsapid.ManifestType(v.(string))
	}

	if v, ok := data["os"]; ok {
		manifest.Os = v.(string)
	}

	if v, ok := data["requirements"]; ok {
		manifest.Requirements = dsapid.Table(v.(map[string]interface{}))
	}

	if v, ok := data["tags"]; ok {
		manifest.Tags = dsapid.Table(v.(map[string]interface{}))
	}

	if v, ok := data["users"]; ok {
		for _, u := range v.([]interface{}) {
			manifest.Users = append(manifest.Users, dsapid.Table(u.(map[string]interface{})))
		}
	}

	if v, ok := data["public"]; ok {
		manifest.Public = v.(bool)
	} else {
		manifest.Public = true
	}

	if v, ok := data["disabled"]; ok {
		manifest.Disabled = v.(bool)
	} else {
		manifest.Disabled = false
	}

	if manifest.Disabled {
		manifest.State = dsapid.ManifestStateDisabled
	} else {
		manifest.State = dsapid.ManifestStateActive
	}

	if v, ok := data["homepage"]; ok {
		manifest.Homepage = v.(string)
	}

	if v, ok := data["urn"]; ok {
		manifest.Urn = v.(string)
	} else {
		manifest.Urn = converter.ComputeUrn(manifest)
	}

	if manifest.Type == dsapid.ManifestTypeZvol {
		if v, ok := data["nic_driver"]; ok {
			manifest.Options["nic_driver"] = v.(string)
		}

		if v, ok := data["disk_driver"]; ok {
			manifest.Options["disk_driver"] = v.(string)
		}

		if v, ok := data["cpu_type"]; ok {
			manifest.Options["cpu_type"] = v.(string)
		}

		if v, ok := data["image_size"]; ok {
			manifest.Options["image_size"] = converter.DecodeToInt64(v)
		}
	}

	if v, ok := data["published_at"]; ok {
		if dt, err := converter.ParseDateTime(v.(string)); err == nil {
			manifest.PublishedAt = dt
		}
	}

	if v, ok := data["created_at"]; ok {
		if dt, err := converter.ParseDateTime(v.(string)); err == nil {
			manifest.CreatedAt = dt
		} else {
			manifest.CreatedAt = manifest.PublishedAt
		}
	} else {
		manifest.CreatedAt = manifest.PublishedAt
	}

	decode_file := func(data dsapid.Table) (file dsapid.ManifestFileResource) {
		file.Compression = dsapid.CompressionTypeNone

		if v, ok := data["size"]; ok {
			file.Size = converter.DecodeToInt64(v)
		}

		if v, ok := data["path"]; ok {
			file.Path = v.(string)
		}

		if v, ok := data["md5"]; ok {
			file.Md5 = v.(string)
		}

		if v, ok := data["sha1"]; ok {
			file.Sha1 = v.(string)
		}

		if ext := path.Ext(file.Path); ext != "" {
			if v, ok := dsapid.CompressionExtensionMap[ext[1:]]; ok {
				file.Compression = v
			}
		}

		return file
	}

	if v, ok := data["files"]; ok {
		for _, u := range v.([]interface{}) {
			manifest.Files = append(manifest.Files, decode_file(dsapid.Table(u.(map[string]interface{}))))
		}
	}

	// internal extra stuff which might only be accessible if we sync from another instance of ourselve
	if v, ok := data["builder_info"]; ok {
		manifest.BuilderInfo = dsapid.Table(v.(map[string]interface{}))
	}

	if v, ok := data["metadata_info"]; ok {
		for _, m := range v.([]interface{}) {
			manifest.MetadataInfo = append(manifest.MetadataInfo, dsapid.Table(m.(map[string]interface{})))
		}
	}

	return manifest
}
