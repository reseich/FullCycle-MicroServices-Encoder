package services_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/reseich/FullCycle-MicroServices-Encoder/application/services"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		logrus.Fatalf("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("encoder-fullcycle-storage")

	require.Nil(t, err)

	err = videoService.Fragment()

	require.Nil(t, err)

	err = videoService.Encode()

	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "encoder-fullcycle-storage"
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(50, doneUpload)

	result := <-doneUpload

	require.Equal(t, result, "Upload Completed")

	err = videoService.Finish()

	require.Nil(t, err)
}
