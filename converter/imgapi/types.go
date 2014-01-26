package imgapi

import (
	"github.com/MerlinDMC/dsapid"
	"time"
)

const (
	CurrentManifestVersion uint = 2
)

type imgapiManifest struct {
	V            uint                 `json:"v"`
	Uuid         string               `json:"uuid"`
	Name         string               `json:"name"`
	Version      string               `json:"version"`
	Description  string               `json:"description"`
	Os           string               `json:"os"`
	Type         dsapid.ManifestType  `json:"type"`
	Urn          string               `json:"urn,omitempty"`
	State        dsapid.ManifestState `json:"state"`
	Disabled     bool                 `json:"disabled"`
	Public       bool                 `json:"public"`
	PublishedAt  time.Time            `json:"published_at"`
	Owner        string               `json:"owner"`
	Requirements dsapid.Table         `json:"requirements,omitempty"`
	Users        []dsapid.Table       `json:"users,omitempty"`
	Tags         dsapid.Table         `json:"tags,omitempty"`
	CpuType      string               `json:"cpu_type,omitempty"`
	ImageSize    int64                `json:"image_size,omitempty"`
	NicDriver    string               `json:"nic_driver,omitempty"`
	DiskDriver   string               `json:"disk_driver,omitempty"`
	Files        []imgapiManifestFile `json:"files"`
}

type imgapiManifestFile struct {
	Md5         string                 `json:"md5,omitempty"`
	Sha1        string                 `json:"sha1,omitempty"`
	Size        int64                  `json:"size"`
	Compression dsapid.CompressionType `json:"compression"`
}
