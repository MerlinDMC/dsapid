package dsapid

import (
	"time"
)

type Table map[string]interface{}

type UserType string
type UserRoleName string
type SyncType string
type SyncProvider string
type CompressionType string
type ManifestState string
type ManifestType string

type UserResource struct {
	Uuid     string         `json:"uuid"`
	Name     string         `json:"name"`
	Password string         `json:"password,omitempty"`
	Email    string         `json:"email,omitempty"`
	Token    string         `json:"token,omitempty"`
	Type     UserType       `json:"type,omitempty"`
	Provider SyncProvider   `json:"provider,omitempty"`
	Roles    []UserRoleName `json:"roles,omitempty"`
}

func (me *UserResource) GetId() string {
	return me.Uuid
}

func (me *UserResource) GetName() string {
	return me.Name
}

func (me *UserResource) HasRoles(roles ...UserRoleName) bool {
	var matched_count int = 0

nextRole:
	for _, role := range roles {
		for _, r := range me.Roles {
			if role == r {
				matched_count++
				continue nextRole
			}
		}
	}

	return len(roles) == matched_count
}

func (me *UserResource) IsGuest() bool {
	return me.Name == DefaultUserGuestName
}

func (me *UserResource) GetAuthInfo() interface{} {
	return map[string]string{
		"uuid": me.Uuid,
		"name": me.Name,
	}
}

type SyncSourceResource struct {
	Name       string       `json:"name"`
	Active     bool         `json:"active"`
	Type       SyncType     `json:"type"`
	Provider   SyncProvider `json:"provider"`
	Source     string       `json:"source,omitempty"`
	FileSource string       `json:"file_source,omitempty"`
	Delay      string       `json:"delay"`
	Opts       Table        `json:"opts,omitempty"`
}

type ManifestResource struct {
	Uuid     string       `json:"uuid"`
	Provider SyncProvider `json:"provider,omitempty"`
	Owner    string       `json:"owner,omitempty"`

	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Homepage string `json:"homepage,omitempty"`
	Urn      string `json:"urn,omitempty"`

	State    ManifestState `json:"state"`
	Public   bool          `json:"public"`
	Disabled bool          `json:"disabled"`

	Type ManifestType `json:"type"`
	Os   string       `json:"os"`

	// TODO: acl
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`

	Requirements Table   `json:"requirements"`
	Users        []Table `json:"users,omitempty"`
	Tags         Table   `json:"tags,omitempty"`
	Options      Table   `json:"options,omitempty"`

	MetadataInfo []Table `json:"metadata_info"`
	BuilderInfo  Table   `json:"builder_info"`
	SyncInfo     Table   `json:"sync_info"`
	// TODO: stats

	Files []ManifestFileResource `json:"files"`
}

type ManifestFileResource struct {
	Path        string          `json:"path"`
	Size        int64           `json:"size"`
	Sha1        string          `json:"sha1"`
	Md5         string          `json:"md5"`
	Compression CompressionType `json:"compression"`
}

type ManifestFilter func(*ManifestResource) bool

const (
	DefaultUserGuestName string = "guest"

	UserTypeUser   UserType = "user"
	UserTypeSystem UserType = "system"

	UserRoleAdmin UserRoleName = "admin"
	UserRoleGuest UserRoleName = "guest"

	UserRoleDatasetUpload UserRoleName = "s_dataset.upload"
	UserRoleDatasetManage UserRoleName = "s_dataset.manage"
	UserRoleDatasetAdmin  UserRoleName = "s_dataset.admin"

	DefaultUserUuid string = "a979f956-12cb-4216-bf4c-ae73e6f14dde"
	DefaultUserName string = "sdc"

	DefaultSyncDelay string = "8h"

	SyncTypeDsapi  SyncType = "dsapi"
	SyncTypeImgapi SyncType = "imgapi"

	SyncProviderJoyent    SyncProvider = "joyent"
	SyncProviderEc        SyncProvider = "ec"
	SyncProviderElys      SyncProvider = "elys"
	SyncProviderCommunity SyncProvider = "community"
	SyncProviderTesting   SyncProvider = "testing"

	ManifestStatePending    ManifestState = "pending"
	ManifestStateActive     ManifestState = "active"
	ManifestStateInactive   ManifestState = "unactivated"
	ManifestStateDisabled   ManifestState = "disabled"
	ManifestStateDeprecated ManifestState = "deprecated"
	ManifestStateNuked      ManifestState = "nuked"

	ManifestTypeZone ManifestType = "zone-dataset"
	ManifestTypeZvol ManifestType = "zvol"

	CompressionTypeGzip  CompressionType = "gzip"
	CompressionTypeBzip2 CompressionType = "bzip2"
	CompressionTypeXz    CompressionType = "xz"
	CompressionTypeNone  CompressionType = "none"
)

var (
	UserTypeDescription = map[UserType]string{
		UserTypeUser:   "Real user account",
		UserTypeSystem: "System account",
	}

	SyncTypeDescription = map[SyncType]string{
		SyncTypeDsapi:  "DSAPI sync source",
		SyncTypeImgapi: "IMGAPI sync source",
	}

	SyncProviderDescription = map[SyncProvider]string{
		SyncProviderJoyent:    "Joyent",
		SyncProviderEc:        "EveryCity",
		SyncProviderCommunity: "Community",
		SyncProviderTesting:   "Testing",
	}

	ManifestStateDescription = map[ManifestState]string{
		ManifestStatePending:    "Pending",
		ManifestStateActive:     "Active",
		ManifestStateInactive:   "Not activated",
		ManifestStateDisabled:   "Disabled",
		ManifestStateDeprecated: "Deprecated",
		ManifestStateNuked:      "Nuked",
	}

	ManifestTypeDescription = map[ManifestType]string{
		ManifestTypeZone: "zone dataset",
		ManifestTypeZvol: "KVM volume",
	}

	CompressionExtensionMap = map[string]CompressionType{
		"gz":   CompressionTypeGzip,
		"bz":   CompressionTypeBzip2,
		"bz2":  CompressionTypeBzip2,
		"xz":   CompressionTypeXz,
		"none": CompressionTypeNone,
	}

	CompressionTypeExtensionMap = map[CompressionType]string{
		CompressionTypeGzip:  ".gz",
		CompressionTypeBzip2: ".bz2",
		CompressionTypeXz:    ".xz",
		CompressionTypeNone:  "",
	}
)
