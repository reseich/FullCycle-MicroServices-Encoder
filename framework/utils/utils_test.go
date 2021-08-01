package utils_test

import (
	"testing"

	"github.com/reseich/FullCycle-MicroServices-Encoder/framework/utils"
	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `{
		"id":"ASDJHAKJSDHSA",
		"file_path":"ASDKLAJSDKLSAD",
		"status": "ASDJASLKDJAKLSD"
		}`

	err := utils.IsJson(json)
	require.Nil(t, err)


	json = `{asdasdasdas}`

	err = utils.IsJson(json)
	require.Error(t, err)

}
