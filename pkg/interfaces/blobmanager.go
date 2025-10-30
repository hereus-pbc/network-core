package interfaces

import "github.com/hereus-pbc/network-core/pkg/types"

type BlobManager interface {
	GetBlobMeta(blobType string, username string, filePath string) (*types.BlobMeta, error)
	ReadBlob(blobType string, username string, filePath string, getAs string) ([]byte, error)
	SaveBlob(blobType string, username string, filePath string, data []byte) (string, error)
	DeleteBlob(blobType string, username string, filePath string) error
	SetBlobPermissions(blobType string, username string, filePath string, targetUser string, permissions []bool) error
	DeleteBlobPermissions(blobType string, username string, filePath string, targetUser string) error
	OverwriteBlob(blobType string, username string, filePath string, writeAs string, data []byte) (string, error)
	ReadBlobPermissions(blobType string, username string, filePath string) (map[string][]bool, error)
}
