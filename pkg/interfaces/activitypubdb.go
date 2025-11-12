package interfaces

import "github.com/hereus-pbc/network-core/pkg/types"

//goland:noinspection ALL
type ActivityPubDB interface {
	GetActor(actorUrl string, actorHandler string) (types.Actor, error)

	// Social Graph Related
	ListFollowers(actor string) []string
	ListFollowing(actor string) []string
	ListFriends(actor string) []string
	AddFollow(who string, whom string, date string, id string) error
	RemoveFollow(who string, whom string) (*[]types.FollowModal, error)

	// Note Related
	DeleteNote(noteId string) error
	FetchNoteByUrl(accessorUrl string, noteUrl string) (types.ActivityPubNote, error)
	FetchNoteById(accessorUrl string, noteId string) (types.ActivityPubNote, error)
	FetchActivityById(accessorUrl string, activityId string) (types.ActivityStream, error)
	ListRepliesToNote(accessorUrl string, noteId string) []types.ActivityPubNote
	EditNote(note types.ActivityPubNote) error
	ListNotes(accessorUrl string, filter map[string]interface{}, packageName string, max int64, latestFirst bool) ([]types.ActivityPubNote, error)

	// Like Related
	RecordLike(who string, noteId string, date string, id string) error
	RemoveLike(who string, noteId string) ([]string, error)
	HasLiked(who string, noteId string) bool
	ListLikes(noteId string) []string

	// Announce Related
	HasAnnounced(who string, what string) bool
	RecordAnnounce(who string, what string, date string, id string, to string, cc []string) error
	RemoveAnnounce(who string, what string) (*[]types.ActivityPubAnnouncement, error)
	ListAnnounces(what string) []string
}
