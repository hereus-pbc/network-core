package interfaces

import "github.com/hereus-pbc/network-core/pkg/types"

type ActivityPubDB interface {
	ListFollowers(actor string) []string
	ListFollowing(actor string) []string
	ListFriends(actor string) []string
	AddFollow(who string, whom string, date string, id string) error
	RemoveFollow(who string, whom string) (*[]types.FollowModal, error)
	DeleteNote(noteId string) error
	LikeNote(who string, noteId string, date string, id string) error
	UnlikeNote(who string, noteId string) ([]string, error)
	FetchNoteByUrl(accessorUrl string, noteUrl string) (types.ActivityPubNote, error)
	FetchNoteById(accessorUrl string, noteId string) (types.ActivityPubNote, error)
	FetchActivityById(accessorUrl string, activityId string) (types.ActivityStream, error)
	GetActor(actorUrl string, actorHandler string) (types.Actor, error)
	ListRepliesToNote(accessorUrl string, noteId string) []types.ActivityPubNote
	HasAnnounced(who string, what string) bool
	UndoAnnounce(who string, what string) (*[]types.ActivityPubAnnouncement, error)
	EditNote(note types.ActivityPubNote) error
	ListNotes(accessorUrl string, filter map[string]interface{}, packageName string, max int64, latestFirst bool) ([]types.ActivityPubNote, error)
}
