package api

import (
	"dts/tabasco/storage/types"
	"encoding/json"
	"net/http"
)

type PutBuildRequest struct {
	build *types.Build
}

func CreatePutBuildRequest(r *http.Request) (PutBuildRequest, error) {
	var build types.Build
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		return PutBuildRequest{}, err
	}
	return PutBuildRequest{build: &build}, nil
}
