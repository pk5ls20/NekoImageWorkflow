// Package uuid rewrite https://github.com/hv0905/NekoImageGallery/blob/master/app/util/generate_uuid.py
package uuid

import (
	"crypto/sha1"
	"github.com/google/uuid"
	"github.com/pk5ls20/NekoImageWorkflow/common/log"
	"os"
)

const namespaceStr = "github.com/hv0905/NekoImageGallery"

// generateUUID is a private function that generates a UUID from a byte slice.
func generateUUID(data []byte) uuid.UUID {
	namespaceUUID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(namespaceStr))
	dataHash := sha1.New()
	dataHash.Write(data)
	return uuid.NewSHA1(namespaceUUID, dataHash.Sum(nil))
}

// GenerateFileUUID generates a UUID based on file content.
func GenerateFileUUID(filePath string) (uuid.UUID, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return uuid.Nil, log.ErrorWrap(err)
	}
	return generateUUID(data), nil
}

// GenerateStrUUID generates a UUID based on the provided string.
func GenerateStrUUID(input string) uuid.UUID {
	return generateUUID([]byte(input))
}

// GenerateByteSliceUUID generates a UUID based on a byte slice.
func GenerateByteSliceUUID(data []byte) uuid.UUID {
	return generateUUID(data)
}
