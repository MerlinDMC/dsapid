package handler

import (
	"encoding/json"
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/MerlinDMC/dsapid/storage"
	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/pborman/uuid"
	"io"
	"io/ioutil"
	"net/http"
)

func ApiGetUsers(encoder middleware.OutputEncoder, params martini.Params, users storage.UserStorage, user middleware.User, req *http.Request) (int, []byte) {
	var users_list []*dsapid.UserResource

	for _, u := range users.Dump() {
		users_list = append(users_list, u)
	}

	return http.StatusOK, encoder.MustEncode(users_list)
}

func ApiPutUsers(encoder middleware.OutputEncoder, params martini.Params, users storage.UserStorage, user middleware.User, req *http.Request) (int, []byte) {
	decoder := json.NewDecoder(req.Body)

	for {
		var u dsapid.UserResource

		if err := decoder.Decode(&u); err == io.EOF {
			break
		} else if err != nil {
			req.Body.Close()

			return http.StatusInternalServerError, encoder.MustEncode(dsapid.Table{
				"error": "invalid users stream",
			})
		}

		// skip empty usernames
		if u.Name == "" {
			continue
		}

		if u.Uuid == "" {
			u.Uuid = uuid.New()
		}

		users.Add(u.Uuid, &u)
	}

	return http.StatusOK, encoder.MustEncode(dsapid.Table{
		"ok": "user(s) added",
	})
}

func ApiUpdateUser(encoder middleware.OutputEncoder, params martini.Params, users storage.UserStorage, user middleware.User, req *http.Request) (int, []byte) {
	action := req.URL.Query().Get("action")
	data, _ := ioutil.ReadAll(req.Body)

	if u, ok := users.GetOK(params["id"]); ok {
		switch action {
		case "set_token":
			log.WithFields(log.Fields{
				"user_uuid": user.GetId(),
				"user_name": user.GetName(),
				"token":     data,
			}).Info("setting token for user")

			u.Token = string(data)

			users.Update(u.Uuid, u)
			break
		case "add_role":
			log.WithFields(log.Fields{
				"user_uuid": user.GetId(),
				"user_name": user.GetName(),
				"new_role":  data,
			}).Info("adding role to user")

			role_name := dsapid.UserRoleName(data)

			if !u.HasRoles(role_name) {
				u.Roles = append(u.Roles, role_name)
			}

			users.Update(u.Uuid, u)
			break
		case "remove_role":
			log.WithFields(log.Fields{
				"user_uuid": user.GetId(),
				"user_name": user.GetName(),
				"new_role":  data,
			}).Info("removing role from user")

			role_name := dsapid.UserRoleName(data)

			if u.HasRoles(role_name) {
				for i, r := range u.Roles {
					if r == role_name {
						u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)

						break
					}
				}
			}

			users.Update(u.Uuid, u)
			break
		}

		return http.StatusOK, encoder.MustEncode(dsapid.Table{
			"ok": "user updated",
		})
	}

	return http.StatusNotFound, encoder.MustEncode(dsapid.Table{
		"error": "user not found",
	})
}

func ApiDeleteUser(encoder middleware.OutputEncoder, params martini.Params, users storage.UserStorage, user middleware.User, req *http.Request) (int, []byte) {
	if u, ok := users.GetOK(params["id"]); ok == true {
		log.WithFields(log.Fields{
			"user_uuid": u.GetId(),
			"user_name": u.GetName(),
		}).Info("deleting user")

		users.Delete(u.Uuid)

		return http.StatusOK, encoder.MustEncode(dsapid.Table{
			"ok": "user deleted",
		})
	}

	return http.StatusNotFound, encoder.MustEncode(dsapid.Table{
		"error": "user not found",
	})
}
