package sync

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/server/logger"
	"github.com/MerlinDMC/dsapid/storage"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type Syncer interface {
	Init(chan *syncerDownloadJob) error
	Run(chan struct{}) error
}

type SyncManager interface {
	Init() error
	Run() error
	Stop()
	Add(Syncer) error
	NewSyncer(dsapid.SyncSourceResource) error
}

type syncerDownloadJob struct {
	manifest *dsapid.ManifestResource
	files    []*url.URL
}

type syncManager struct {
	ParallelFetches int
	users           storage.UserStorage
	manifests       storage.ManifestStorage

	syncer     []Syncer
	q_download chan *syncerDownloadJob
	s_stop     chan struct{}
}

func NewManager(parallel_fetches int, users storage.UserStorage, manifests storage.ManifestStorage) SyncManager {
	manager := &syncManager{
		ParallelFetches: parallel_fetches,
		users:           users,
		manifests:       manifests,
	}

	return manager
}

func (me *syncManager) Init() error {
	me.syncer = make([]Syncer, 0)
	me.q_download = make(chan *syncerDownloadJob)

	return nil
}

func (me *syncManager) Run() error {
	if me.s_stop != nil {
		return ErrSyncAlreadyRunning
	}

	me.s_stop = make(chan struct{})

	if me.ParallelFetches < 1 {
		me.ParallelFetches = 1
	}

	for i := 0; i < me.ParallelFetches; i++ {
		go me.processDownloadJobs()
	}

	return nil
}

func (me *syncManager) Stop() {
	if me.s_stop != nil {
		close(me.s_stop)
	}
}

func (me *syncManager) Add(syncer Syncer) error {
	me.syncer = append(me.syncer, syncer)

	syncer.Init(me.q_download)

	return syncer.Run(me.s_stop)
}

func (me *syncManager) NewSyncer(source dsapid.SyncSourceResource) error {
	var syncer Syncer

	switch source.Type {
	case dsapid.SyncTypeDsapi:
		syncer = &dsapiSyncer{
			source:    &source,
			users:     me.users,
			manifests: me.manifests,
		}
		break
	case dsapid.SyncTypeImgapi:
		syncer = &imgapiSyncer{
			source:    &source,
			users:     me.users,
			manifests: me.manifests,
		}
		break
	}

	return me.Add(syncer)
}

func (me *syncManager) processDownloadJobs() {
	for {
		select {
		case <-me.s_stop:
			// received stop signal -> exit processing loop
			logger.Info("received stop signal. exiting processing loop.")
			return
		case job := <-me.q_download:
			if _, ok := me.manifests.GetOK(job.manifest.Uuid); !ok {
				logger.Tracef("need to fetch new image: %s (%s v%s)", job.manifest.Uuid, job.manifest.Name, job.manifest.Version)

				if err := os.MkdirAll(me.manifests.ManifestPath(job.manifest), 0770); err == nil {
					for file_idx, src := range job.files {
						filename := me.manifests.FilePath(job.manifest, &job.manifest.Files[file_idx])

						if err := me.downloadManifestFile(src, filename, &job.manifest.Files[file_idx]); err != nil {
							logger.Errorf("download error: %s", err)
						} else {
							me.manifests.Add(job.manifest.Uuid, job.manifest)
						}
					}
				}
			}

			break
		}
	}
}

func (me *syncManager) downloadManifestFile(src *url.URL, filename string, file *dsapid.ManifestFileResource) (err error) {
	hash_md5 := md5.New()
	hash_sha1 := sha1.New()

	wget := exec.Command("wget", "-c", "--no-check-certificate", src.String(), "-O", filename)

	logger.Debugf("running wget: %s", strings.Join(wget.Args, " "))

	if err := wget.Run(); err == nil {
		if file_in, err := os.OpenFile(filename, os.O_RDONLY, 0660); err == nil {
			defer file_in.Close()

			writer := io.MultiWriter(hash_md5, hash_sha1)

			if _, err := io.Copy(writer, file_in); err == nil {
				md5_sum := hex.EncodeToString(hash_md5.Sum(nil))
				sha1_sum := hex.EncodeToString(hash_sha1.Sum(nil))

				if file.Md5 != "" && file.Md5 != md5_sum {
					logger.Warnf("checksum mismatch: got %s expected %s", md5_sum, file.Md5)
					return ErrChecksumNotMatching
				}

				if file.Sha1 != "" && file.Sha1 != sha1_sum {
					logger.Warnf("checksum mismatch: got %s expected %s", sha1_sum, file.Sha1)
					return ErrChecksumNotMatching
				}

				file.Md5 = md5_sum
				file.Sha1 = sha1_sum
			} else {
				logger.Errorf("%s", err)
				return err
			}
		}
	} else {
		logger.Errorf("%s", err)
		return err
	}

	return
}
