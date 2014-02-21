package main

import (
	"github.com/MerlinDMC/dsapid"
	"github.com/MerlinDMC/dsapid/server/handler"
	"github.com/MerlinDMC/dsapid/server/middleware"
	"github.com/codegangsta/martini"
)

func registerRoutes(router martini.Router) {
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
	router.Get("/api/datasets", middleware.AllowCORS(), handler.ApiDatasetsList)
	router.Get("/api/datasets/:id", middleware.AllowCORS(), handler.ApiDatasetsDetail)
	router.Get("/api/export/:id", handler.ApiDatasetExport)

	// private api - update
	router.Post("/api/reload/datasets", middleware.RequireRoles(dsapid.UserRoleDatasetAdmin), handler.ApiPostReloadDatasets)
	router.Post("/api/datasets/:id", middleware.RequireRoles(dsapid.UserRoleDatasetManage), handler.ApiPostDatasetUpdate)

	// private api - upload
	router.Post("/api/upload", middleware.RequireRoles(dsapid.UserRoleDatasetUpload), handler.ApiPostFileUpload)
}
