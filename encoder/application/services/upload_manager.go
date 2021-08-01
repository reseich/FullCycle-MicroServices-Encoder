package services

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectpath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectpath, os.Getenv("LOCAL_STORAGE_PATH")+"/")
	f, err := os.Open(objectpath)

	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	// wc.ACL = []storage.ACLRule{{Entity: storage.AllAuthenticatedUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) LoadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.LoadPaths()

	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()

	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil

}

func (vu *VideoUpload) uploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {
	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)
		logrus.Printf("File %s ", vu.Paths[x])
		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			logrus.Printf("erro during the upload: %s, Error:%v", vu.Paths[x], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}
	returnChan <- "Upload Completed"
}
