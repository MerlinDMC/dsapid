package main

import (
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/server/handler"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/go-martini/martini"
)

func registerRoutes(router martini.Router, config Config) {
	// common
	router.Get("/ping", handler.CommonPing)
	router.Get("/status", handler.CommonStatus)

	// dsapi
	router.Get("/datasets", middleware.AllowCORS(), handler.DsapiList)
	router.Get("/datasets/:id", middleware.AllowCORS(), handler.DsapiDetail)
	router.Get("/datasets/:id/:path", handler.DsapiFile)

	// imgapi
	router.Get("/images", middleware.AllowCORS(), handler.ImgapiList)
	router.Get("/images/:id", middleware.AllowCORS(), handler.ImgapiDetail)
	router.Get("/images/:id/file", handler.ImgapiFile)
	router.Get("/images/:id/file:file_idx", handler.ImgapiFile)

	// public api
	router.Group("/api", func(router martini.Router) {
		router.Get("/datasets", middleware.AllowCORS(), handler.ApiDatasetsList)
		router.Get("/datasets/:id", middleware.AllowCORS(), handler.ApiDatasetsDetail)
		router.Get("/export/:id", handler.ApiDatasetExport)
	}, middleware.Throttle(config.Throttle.Api.ToQuota()))

	// private api - update
	router.Group("/api", func(router martini.Router) {
		router.Post("/reload/datasets", handler.ApiPostReloadDatasets)
		router.Post("/datasets/:id", handler.ApiPostDatasetUpdate)
	}, middleware.RequireRoles(dsapid.UserRoleDatasetAdmin))

	// private api - users
	router.Group("/api", func(router martini.Router) {
		router.Get("/users", handler.ApiGetUsers)
		router.Put("/users", handler.ApiPutUsers)
		router.Post("/users/:id", handler.ApiUpdateUser)
		router.Delete("/users/:id", handler.ApiDeleteUser)
	}, middleware.RequireAdmin())

	// private api - upload
	router.Post("/api/upload", middleware.RequireRoles(dsapid.UserRoleDatasetUpload), handler.ApiPostFileUpload)
}
