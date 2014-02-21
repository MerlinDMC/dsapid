package storage

import (
	"github.com/MerlinDMC/dsapid"
	"strings"
)

func FilterManifestEnabled() ManifestFilter {
	return func(manifest *dsapid.ManifestResource) bool {
		return manifest.State == dsapid.ManifestStateActive || manifest.State == dsapid.ManifestStateDeprecated
	}
}

func FilterManifestPublic(value bool) ManifestFilter {
	return func(manifest *dsapid.ManifestResource) bool {
		return manifest.Public == value
	}
}

func FilterManifestForUser(uuid string) ManifestFilter {
	return func(manifest *dsapid.ManifestResource) bool {
		return manifest.Public == true || manifest.Owner == uuid
	}
}

func FilterManifestName(value string) ManifestFilter {
	return func(manifest *dsapid.ManifestResource) bool {
		if strings.HasPrefix(value, "~") {
			return strings.Contains(manifest.Name, value[1:])
		} else {
			return strings.HasPrefix(manifest.Name, value)
		}
	}
}

func FilterManifestVersion(value string) ManifestFilter {
	return func(manifest *dsapid.ManifestResource) bool {
		return strings.HasPrefix(manifest.Version, value)
	}
}

func FilterManifestOs(value string) ManifestFilter {
	return func(manifest *dsapid.ManifestResource) bool {
		return strings.HasPrefix(manifest.Os, value)
	}
}

func FilterManifestUuid(value string) ManifestFilter {
	return func(manifest *dsapid.ManifestResource) bool {
		return strings.HasPrefix(manifest.Uuid, value)
	}
}
