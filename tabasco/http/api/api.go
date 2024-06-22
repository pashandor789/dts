package api

import (
	pkghttp "dts/pkg/http"
	"dts/pkg/log"
	"dts/tabasco/storage"
	"github.com/go-chi/chi"
	"net/http"
)

type PublicHandler struct {
	logger  log.Logger
	storage *storage.Storage
}

func NewPublicHandler(logger log.Logger, storage *storage.Storage) PublicHandler {
	return PublicHandler{
		logger:  logger,
		storage: storage,
	}
}

func (h *PublicHandler) WithAPIHandlers() pkghttp.RouterOption {
	return func(r chi.Router) {
		pkghttp.AddHandler(r.Get, "/builds", h.getBuilds)

		pkghttp.AddHandler(r.Put, "/build", h.putBuild)
		pkghttp.AddHandler(r.Put, "/tests", h.putTests)
	}
}

// @Summary      Get builds
// @Description  Retrieves a list of builds.
// @Produce      json
// @Success      200  {array}   types.Build   "List of builds"
// @Failure      400  {object}  http.Error    "Bad request"
// @Router       /builds [get]
func (h *PublicHandler) getBuilds(r *http.Request) pkghttp.Response {
	builds, err := h.storage.GetBuilds()
	if err != nil {
		return pkghttp.BadRequest(err)
	}
	return pkghttp.OK(builds)
}

// @Summary 	Put build
// @Description Put build.
// @Accept  	json
// @Produce  	json
// @Param   	build  body   types.Build    true
// @Success 	200 {object}  http.Success	 "ok"
// @Failure 	400 {object}  http.Error	 "bad request"
// @Router 		/build [put]
func (h *PublicHandler) putBuild(r *http.Request) pkghttp.Response {
	req, err := CreatePutBuildRequest(r)
	if err != nil {
		return pkghttp.BadRequest(err)
	}
	err = h.storage.PutBuild(req.build)
	if err != nil {
		return pkghttp.BadRequest(err)
	}
	return pkghttp.OK(nil)
}

func (h *PublicHandler) putTests(r *http.Request) pkghttp.Response {

}
