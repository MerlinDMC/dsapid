package dsapi

import (
	"github.com/MerlinDMC/dsapid"
	"time"
)

type dsapiManifest struct {
	Uuid         string              `json:"uuid"`
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Description  string              `json:"description"`
	Os           string              `json:"os"`
	Type         dsapid.ManifestType `json:"type"`
	Homepage     string              `json:"homepage,omitempty"`
	Urn          string              `json:"urn,omitempty"`
	PublishedAt  time.Time           `json:"published_at"`
	CreatedAt    time.Time           `json:"created_at"`
	CreatorUuid  string              `json:"creator_uuid"`
	CreatorName  string              `json:"creator_name,omitempty"`
	VendorUuid   string              `json:"vendor_uuid"`
	Requirements dsapid.Table        `json:"requirements,omitempty"`
	Users        []dsapid.Table      `json:"users,omitempty"`
	Tags         dsapid.Table        `json:"tags,omitempty"`
	CpuType      string              `json:"cpu_type,omitempty"`
	ImageSize    int64               `json:"image_size,omitempty"`
	NicDriver    string              `json:"nic_driver,omitempty"`
	DiskDriver   string              `json:"disk_driver,omitempty"`
	BuilderInfo  dsapid.Table        `json:"builder_info,omitempty"`
	MetadataInfo []dsapid.Table      `json:"metadata_info,omitempty"`
	Files        []dsapiManifestFile `json:"files"`
}

type dsapiManifestFile struct {
	Url  string `json:"url"`
	Path string `json:"path"`
	Md5  string `json:"md5,omitempty"`
	Sha1 string `json:"sha1,omitempty"`
	Size int64  `json:"size"`
}
