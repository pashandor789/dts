package api

import (
	pkghttp "dts/pkg/http"
	"dts/pkg/log"
	"dts/tabasco/storage"
	"fmt"
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
		pkghttp.AddHandler(r.Get, "/tests/{id}", h.getTests)

		pkghttp.AddHandler(r.Put, "/build", h.putBuild)
		pkghttp.AddHandler(r.Put, "/tests", h.putTests)
	}
}

// @Summary      Retrieves a list of all builds
// @Description  Retrieves a list of all builds.
// @Produce      json
// @Success      200  {array}   types.Build   "List of builds"
// @Failure      400  {object}  http.Error    "Bad Request"
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
// @Param   	build  body   types.Build    required  "Build"
// @Success 	200 {object}  http.Success	 "OK"
// @Failure 	400 {object}  http.Error	 "Bad Request"
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

// @Summary      Retrieve tests by task ID
// @Description  Get tests by the task ID provided as an URL parameter
// @Produce      json
// @Param		 id path string true "Task ID"
// @Success      200  {array}   types.Test   "List of tests"
// @Failure      400  {object}  http.Error    "Bad Request"
// @Failure      404  {object}  http.Error    "Not Found"
// @Router       /tests/{id} [get]
func (h *PublicHandler) getTests(r *http.Request) pkghttp.Response {
	taskID := chi.URLParam(r, "id")
	tests, err := h.storage.GetTests(taskID)
	if err != nil {
		return pkghttp.BadRequest(err)
	}
	if len(tests) == 0 {
		return pkghttp.NotFound(fmt.Errorf("no such tests with taskId: %s", taskID))
	}
	return pkghttp.OK(tests)
}

// @Summary 	Put tests.
// @Description Put tests with multipart/form-data : meta.json, {i}_input, {i}_output or meta.json, tests.zip
// @Accept 		multipart/form-data
// @Produce 	json
// @Param 		tests.zip formData file false "tests.zip with {i}_input, {i}_output"
// @Param 		meta.json formData file true "Meta file"
// @Success 	200 {object}  http.Success	 "ok"
// @Failure 	400 {object}  http.Error	 "bad request"
// @Router 		/tests [put]
func (h *PublicHandler) putTests(r *http.Request) pkghttp.Response {
	req, err := CreatePutTestsRequest(r)
	if err != nil {
		return pkghttp.BadRequest(err)
	}
	err = h.storage.PutTests(req.tests)
	if err != nil {
		pkghttp.BadRequest(err)
	}
	return pkghttp.OK(nil)
}
