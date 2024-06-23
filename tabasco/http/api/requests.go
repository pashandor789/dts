package api

import (
	"archive/zip"
	"bytes"
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

	var getFileContent func(name string) ([]byte, error)
	var testSize int

	testsZipFile, _, err := r.FormFile("tests.zip")
	if err == nil {
		var fileByName map[string][]byte = make(map[string][]byte)

		testsZip, err := io.ReadAll(testsZipFile)
		if err != nil {
			return nil, err
		}

		zipReader, err := zip.NewReader(bytes.NewReader(testsZip), int64(len(testsZip)))
		for _, zipFile := range zipReader.File {
			rc, err := zipFile.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			b, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}
			fileByName[zipFile.Name] = b
		}

		getFileContent = func(name string) ([]byte, error) {
			val, ok := fileByName[fmt.Sprintf("tests/%s", name)]
			if !ok {
				return nil, fmt.Errorf("no such file : %s", name)
			}
			return val, nil
		}
		testSize = len(fileByName) / 2
	} else {
		getFileContent = func(name string) ([]byte, error) {
			f, _, err := r.FormFile(name)
			if err != nil {
				return nil, err
			}

			b, err := io.ReadAll(f)
			return b, nil
		}
		testSize = (len(r.MultipartForm.Value) - 1) / 2
	}

	var putTaskRequests PutTestsRequest
	for i := 1; i <= testSize; i++ {
		inputTest, err := getFileContent(fmt.Sprintf("%d_input", i))
		if err != nil {
			return nil, err
		}
		it := types.Test{Id: uint16(i), Type: types.Input, TaskId: metaFile.TaskId, Data: inputTest}
		putTaskRequests.tests = append(putTaskRequests.tests, it)

		outputTest, err := getFileContent(fmt.Sprintf("%d_output", i))
		if err != nil {
			return nil, err
		}
		ot := types.Test{Id: uint16(i), Type: types.Output, TaskId: metaFile.TaskId, Data: outputTest}
		putTaskRequests.tests = append(putTaskRequests.tests, ot)
	}

	return &putTaskRequests, err
}
