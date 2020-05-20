package keys

import (
	"flag"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE(jalextowle): We must ignore this flag to prevent the flag package from
// panicking when this flag is provided to `wasmbrowsertest` in the browser tests.
func init() {
	_ = flag.String("initFile", "", "")
}

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
