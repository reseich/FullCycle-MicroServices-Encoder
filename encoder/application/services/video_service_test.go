package services_test

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/reseich/FullCycle-MicroServices-Encoder/application/repositories"
	"github.com/reseich/FullCycle-MicroServices-Encoder/application/services"
	"github.com/reseich/FullCycle-MicroServices-Encoder/domain"
	"github.com/reseich/FullCycle-MicroServices-Encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		logrus.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)
	return video, repo
}

func TestVideoServiceDownload(t *testing.T) {
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

	err = videoService.Finish()

	require.Nil(t, err)

}
