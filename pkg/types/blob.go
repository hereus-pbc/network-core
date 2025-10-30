package types

type BlobMeta struct {
	BlobType    string            `bson:"blob_type"`
	Username    string            `bson:"username"`
	FilePath    string            `bson:"file_path"`
	BlobId      string            `bson:"blob_id"`
	Permissions map[string][]bool `bson:"permissions"` // {"username@example.com": [true, false, false, false]} // read, write, comment, admin
}
