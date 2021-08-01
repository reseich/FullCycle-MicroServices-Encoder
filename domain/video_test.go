package domain_test

import (
	"testing"
	"time"

	"github.com/reseich/FullCycle-MicroServices-Encoder/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()
	require.Error(t, err)
}

func TestVideoIdIsNotAUuid(t *testing.T) {
	video := domain.NewVideo()
	video.ID = "abc"
	video.CreatedAt = time.Now()
	video.FilePath = "path"
	video.ResourceID = "test"
	err := video.Validate()
	require.Error(t, err)
}

func TestVideoValidator(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.CreatedAt = time.Now()
	video.FilePath = "path"
	video.ResourceID = "test"
	err := video.Validate()
	require.NoError(t, err)
}
