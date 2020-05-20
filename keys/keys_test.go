package keys

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndGetKey(t *testing.T) {
	path := "/tmp/keys/" + uuid.New().String()
	generatedKey, err := GenerateAndSavePrivateKey(path)
	require.NoError(t, err)
	gotKey, err := GetPrivateKeyFromPath(path)
	require.NoError(t, err)
	assert.Equal(t, generatedKey, gotKey)
}

func TestGetKeyNotExists(t *testing.T) {
	nonExistentPath := "/tmp/keys/" + uuid.New().String()
	_, err := GetPrivateKeyFromPath(nonExistentPath)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err), "error should be a NotExist error, but got: (%T) %s", err, err)
}
