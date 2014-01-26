package storage

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/MerlinDMC/dsapid"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

const (
	defaultManifestFilename string = "manifest.json"
)

type ManifestStorage interface {
	Add(string, *dsapid.ManifestResource)
	Get(string) *dsapid.ManifestResource
	GetOK(string) (*dsapid.ManifestResource, bool)
	List() chan *dsapid.ManifestResource
	Filter(...ManifestFilter) chan *dsapid.ManifestResource
	ManifestPath(*dsapid.ManifestResource) string
	FilePath(*dsapid.ManifestResource, *dsapid.ManifestFileResource) string
}

type filesystemManifestStorage struct {
	basedir string

	manifests map[string]*dsapid.ManifestResource
	byDate    []*dsapid.ManifestResource
}

type ManifestFilter func(*dsapid.ManifestResource) bool

type ManifestsByPublishedAt []*dsapid.ManifestResource

func (t ManifestsByPublishedAt) Len() int      { return len(t) }
func (t ManifestsByPublishedAt) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t ManifestsByPublishedAt) Less(i, j int) bool {
	return t[i].PublishedAt.Unix() > t[j].PublishedAt.Unix()
}

func NewManifestStorage(basedir string) ManifestStorage {
	storage := new(filesystemManifestStorage)

	storage.basedir = basedir

	storage.manifests = make(map[string]*dsapid.ManifestResource)
	storage.byDate = make([]*dsapid.ManifestResource, 0)

	storage.load()

	return storage
}

func (me *filesystemManifestStorage) Add(id string, manifest *dsapid.ManifestResource) {
	os.MkdirAll(path.Join(me.basedir, id), 0770)

	if data, err := json.MarshalIndent(manifest, "", "  "); err == nil {
		err = ioutil.WriteFile(path.Join(me.basedir, id, defaultManifestFilename), data, 0666)
	}

	me.add(id, manifest)
}

func (me *filesystemManifestStorage) Get(id string) *dsapid.ManifestResource {
	return me.manifests[id]
}

func (me *filesystemManifestStorage) GetOK(id string) (*dsapid.ManifestResource, bool) {
	v, ok := me.manifests[id]

	return v, ok
}

func (me *filesystemManifestStorage) List() (c chan *dsapid.ManifestResource) {
	c = make(chan *dsapid.ManifestResource)

	go func() {
		for _, item := range me.byDate {
			c <- item
		}

		close(c)
	}()

	return
}

func (me *filesystemManifestStorage) Filter(flist ...ManifestFilter) (c chan *dsapid.ManifestResource) {
	c = make(chan *dsapid.ManifestResource)

	go func() {
	nextItem:
		for item := range me.List() {
			for _, f := range flist {
				if !f(item) {
					continue nextItem
				}
			}

			c <- item
		}

		close(c)
	}()

	return
}

func (me *filesystemManifestStorage) ManifestPath(manifest *dsapid.ManifestResource) string {
	return path.Join(me.basedir, manifest.Uuid)
}

func (me *filesystemManifestStorage) FilePath(manifest *dsapid.ManifestResource, file *dsapid.ManifestFileResource) string {
	return path.Join(me.basedir, manifest.Uuid, file.Path)
}

func (me *filesystemManifestStorage) add(id string, manifest *dsapid.ManifestResource) {
	me.manifests[id] = manifest
	me.byDate = append(me.byDate, manifest)

	sort.Sort(ManifestsByPublishedAt(me.byDate))
}

func (me *filesystemManifestStorage) load() {
	if items, err := ioutil.ReadDir(me.basedir); err == nil {
		for _, item := range items {
			if item.IsDir() {
				manifestFilename := path.Join(me.basedir, item.Name(), defaultManifestFilename)

				if _, err := os.Stat(manifestFilename); err == nil {
					if data, err := ioutil.ReadFile(manifestFilename); err == nil {
						id, manifest := uuid.Parse(item.Name()), dsapid.ManifestResource{}

						json.Unmarshal(data, &manifest)

						if id.String() == manifest.Uuid {
							me.add(manifest.Uuid, &manifest)
						}
					}
				}
			}
		}
	}
}
