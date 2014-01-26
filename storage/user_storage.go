package storage

import (
	"encoding/json"
	"github.com/MerlinDMC/dsapid"
	"io/ioutil"
	"os"
)

type UserStorage interface {
	Save() error
	Add(string, *dsapid.UserResource)
	Delete(string)
	EnsureExists(string, string) *dsapid.UserResource
	Get(string) *dsapid.UserResource
	GetOK(string) (*dsapid.UserResource, bool)
	FindByName(string) (*dsapid.UserResource, error)
	FindByEmail(string) (*dsapid.UserResource, error)
	FindByToken(string) (*dsapid.UserResource, error)
	GuestUser() *dsapid.UserResource
}

func NewUserStorage(filename string) UserStorage {
	store := new(jsonUserStorage)

	store.filename = filename
	store.loaded = false
	store.users = make(map[string]*dsapid.UserResource, 0)
	store.map_name_id = make(map[string]string)
	store.map_email_id = make(map[string]string)
	store.map_token_id = make(map[string]string)

	return store
}

type jsonUserStorage struct {
	filename string
	loaded   bool
	users    map[string]*dsapid.UserResource

	map_name_id  map[string]string
	map_email_id map[string]string
	map_token_id map[string]string
}

func (me *jsonUserStorage) Save() error {
	return me.save()
}

func (me *jsonUserStorage) Add(id string, user *dsapid.UserResource) {
	me.add(id, *user)
	me.save()
}

func (me *jsonUserStorage) Delete(id string) {
	me.delete(id)
	me.save()
}

func (me *jsonUserStorage) EnsureExists(id, name string) *dsapid.UserResource {
	if v, ok := me.users[id]; ok {
		return v
	}

	user := new(dsapid.UserResource)
	user.Uuid = id
	user.Name = name

	me.Add(id, user)

	return user
}

func (me *jsonUserStorage) Get(id string) *dsapid.UserResource {
	me.load()

	if v, ok := me.users[id]; ok {
		return v
	}

	return nil
}

func (me *jsonUserStorage) GetOK(id string) (*dsapid.UserResource, bool) {
	me.load()

	v, ok := me.users[id]

	return v, ok
}

func (me *jsonUserStorage) FindByName(name string) (*dsapid.UserResource, error) {
	me.load()

	if v, ok := me.map_name_id[name]; ok {
		return me.users[v], nil
	}

	return nil, ErrStorageItemNotFound
}

func (me *jsonUserStorage) FindByEmail(email string) (*dsapid.UserResource, error) {
	me.load()

	if v, ok := me.map_email_id[email]; ok {
		return me.users[v], nil
	}

	return nil, ErrStorageItemNotFound
}

func (me *jsonUserStorage) FindByToken(token string) (*dsapid.UserResource, error) {
	me.load()

	if v, ok := me.map_token_id[token]; ok {
		return me.users[v], nil
	}

	return nil, ErrStorageItemNotFound
}

func (me *jsonUserStorage) GuestUser() *dsapid.UserResource {
	return &dsapid.UserResource{
		Name:  dsapid.DefaultUserGuestName,
		Type:  dsapid.UserTypeSystem,
		Roles: []dsapid.UserRoleName{dsapid.UserRoleGuest},
	}
}

func (me *jsonUserStorage) add(id string, user dsapid.UserResource) {
	me.users[id] = &user

	if _, ok := me.map_name_id[user.Name]; !ok && user.Name != "" {
		me.map_name_id[user.Name] = id
	}

	if _, ok := me.map_email_id[user.Email]; !ok && user.Email != "" {
		me.map_email_id[user.Email] = id
	}

	if _, ok := me.map_token_id[user.Token]; !ok && user.Token != "" {
		me.map_token_id[user.Token] = id
	}
}

func (me *jsonUserStorage) delete(id string) {
	if v, ok := me.users[id]; ok {
		delete(me.users, id)

		if _, ok := me.map_name_id[v.Name]; ok {
			delete(me.map_name_id, v.Name)
		}

		if _, ok := me.map_email_id[v.Email]; ok {
			delete(me.map_email_id, v.Email)
		}

		if _, ok := me.map_token_id[v.Token]; ok {
			delete(me.map_token_id, v.Token)
		}
	}
}

func (me *jsonUserStorage) load() error {
	if me.loaded {
		return nil
	}

	if _, err := os.Stat(me.filename); err == nil {
		if data, err := ioutil.ReadFile(me.filename); err == nil {
			var users []dsapid.UserResource

			if err = json.Unmarshal(data, &users); err == nil {
				for _, u := range users {
					me.add(u.Uuid, u)
				}
			} else {
				return ErrStorageFileInvalid
			}
		} else {
			return ErrStorageFileNotReadable
		}
	} else {
		return ErrStorageFileNotFound
	}

	return nil
}

func (me *jsonUserStorage) save() (err error) {
	if me.filename == "" {
		return
	}

	var users []*dsapid.UserResource = make([]*dsapid.UserResource, 0)

	for _, u := range me.users {
		users = append(users, u)
	}

	if data, err := json.MarshalIndent(users, "", "  "); err == nil {
		err = ioutil.WriteFile(me.filename, data, 0660)
	}

	return err
}
