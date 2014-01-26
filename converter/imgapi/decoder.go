package imgapi

import (
	"fmt"
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/converter"
	"github.com/MerlinDMC/dsapid/storage"
)

type imgadmDecoder struct {
	provider dsapid.SyncProvider
	users    storage.UserStorage
}

func NewDecoder(provider dsapid.SyncProvider, users storage.UserStorage) (decoder *imgadmDecoder, err error) {
	decoder = new(imgadmDecoder)

	if provider != "" {
		decoder.provider = provider
	} else {
		decoder.provider = dsapid.SyncProviderCommunity
	}

	decoder.users = users

	return decoder, nil
}

func (me *imgadmDecoder) Decode(data dsapid.Table) *dsapid.ManifestResource {
	var manifest *dsapid.ManifestResource

	manifest = new(dsapid.ManifestResource)

	manifest.Options = make(dsapid.Table)
	manifest.SyncInfo = make(dsapid.Table)
	manifest.BuilderInfo = make(dsapid.Table)
	manifest.MetadataInfo = make([]dsapid.Table, 0)

	var creator_uuid string = dsapid.DefaultUserUuid
	var creator_name string = dsapid.DefaultUserName

	if v, ok := data["owner"]; ok && v.(string) != "" {
		creator_uuid = v.(string)

		if u, ok := me.users.GetOK(creator_uuid); ok {
			creator_name = u.Name
			// } else {
			// 	logger.Warn("owner unknown will use default name '%s'", creator_name)
		}
		// } else {
		// 	logger.Warn("manifest has no owner")
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

	if v, ok := data["state"]; ok {
		manifest.State = dsapid.ManifestState(v.(string))
	} else {
		if manifest.Disabled {
			manifest.State = dsapid.ManifestStateDisabled
		} else {
			manifest.State = dsapid.ManifestStateActive
		}
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

		if v, ok := data["compression"]; ok {
			file.Compression = dsapid.CompressionType(v.(string))
		}

		if v, ok := data["size"]; ok {
			file.Size = converter.DecodeToInt64(v)
		}

		if v, ok := data["md5"]; ok {
			file.Md5 = v.(string)
		}

		if v, ok := data["sha1"]; ok {
			file.Sha1 = v.(string)
		}

		// reconstruct path information
		if ext, ok := dsapid.CompressionTypeExtensionMap[file.Compression]; ok {
			var fullName string
			if file_idx := len(manifest.Files); file_idx > 0 {
				fullName = fmt.Sprintf("%s-%s-%d", manifest.Name, manifest.Version, file_idx)
			} else {
				fullName = fmt.Sprintf("%s-%s", manifest.Name, manifest.Version)
			}

			if manifest.Type == dsapid.ManifestTypeZone {
				file.Path = fmt.Sprintf("%s.zfs%s", fullName, ext)
			} else {
				file.Path = fmt.Sprintf("%s.zvol%s", fullName, ext)
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
