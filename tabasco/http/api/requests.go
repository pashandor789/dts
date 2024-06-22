package api

import (
	"dts/tabasco/storage/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PutBuildRequest struct {
	build *types.Build
}

func CreatePutBuildRequest(r *http.Request) (*PutBuildRequest, error) {
	var build types.Build
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		return nil, err
	}
	return &PutBuildRequest{build: &build}, nil
}

type PutTestsRequest struct {
	tests    []types.Test
	taskMeta types.TaskMeta
}

func CreatePutTestsRequest(r *http.Request) (*PutTestsRequest, error) {
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		return nil, err
	}

	metaFileRaw, _, err := r.FormFile("meta.json")
	if err != nil {
		return nil, err
	}

	var metaFile types.TaskMeta
	if err := json.NewDecoder(metaFileRaw).Decode(&metaFile); err != nil {
		return nil, err
	}

	var putTaskRequests PutTestsRequest
	testSize := (len(r.MultipartForm.Value) - 1) / 2
	for i := 1; i <= testSize; i++ {
		inputFile, _, err := r.FormFile(fmt.Sprintf("%d_input", i))
		if err != nil {
			return nil, err
		}
		inputTest, err := io.ReadAll(inputFile)
		if err != nil {
			return nil, err
		}

		outputFile, _, err := r.FormFile(fmt.Sprintf("%d_output", i))
		if err != nil {
			return nil, err
		}
		outputTest, err := io.ReadAll(outputFile)
		if err != nil {
			return nil, err
		}

		it := types.Test{Id: uint16(i), Type: types.Input, TaskId: metaFile.TaskId, Data: inputTest}
		putTaskRequests.tests = append(putTaskRequests.tests, it)
		ot := types.Test{Id: uint16(i), Type: types.Output, TaskId: metaFile.TaskId, Data: outputTest}
		putTaskRequests.tests = append(putTaskRequests.tests, ot)
	}
}
