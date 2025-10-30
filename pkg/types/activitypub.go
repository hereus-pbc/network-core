package types

type ActivityStream struct {
	LdContext interface{} `json:"@context,omitempty" bson:"@context"`
	Id        string      `json:"id" bson:"id"`
	Type      string      `json:"type" bson:"type"`
	Actor     string      `json:"actor" bson:"actor"`
	Object    interface{} `json:"object" bson:"object"`
	To        interface{} `json:"to" bson:"to,omitempty"`
	Cc        interface{} `json:"cc,omitempty" bson:"cc,omitempty"`
}

type ActorPublicKey struct {
	Id        string `json:"id" bson:"id"`
	Owner     string `json:"owner" bson:"owner"`
	PublicKey string `json:"publicKeyPem" bson:"publicKey"`
}

type ActorIcon struct {
	Type string `json:"type" bson:"type"`
	Url  string `json:"url" bson:"url"`
}

type ActorEndpoints struct {
	SharedInbox string `json:"sharedInbox" bson:"sharedInbox"`
}

type ActorAttachment struct {
	Type  string `json:"type" bson:"type"`
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type Actor struct {
	LdContext                 interface{}       `json:"@context" bson:"@context"`
	Id                        string            `json:"id" bson:"id"`
	Type                      string            `json:"type" bson:"type"`
	Inbox                     string            `json:"inbox" bson:"inbox"`
	Outbox                    string            `json:"outbox" bson:"outbox"`
	Username                  string            `json:"preferredUsername" bson:"preferredUsername"`
	Followers                 string            `json:"followers" bson:"followers"`
	Following                 string            `json:"following" bson:"following"`
	Name                      string            `json:"name" bson:"name"`
	Url                       string            `json:"url" bson:"url"`
	Summary                   string            `json:"summary" bson:"summary"`
	Endpoints                 ActorEndpoints    `json:"endpoints" bson:"endpoints"`
	PublicKey                 ActorPublicKey    `json:"publicKey" bson:"publicKey"`
	Icon                      ActorIcon         `json:"icon" bson:"icon"` // below are Mastodon extras
	ManuallyApprovesFollowers bool              `json:"manuallyApprovesFollowers" bson:"manuallyApprovesFollowers"`
	Discoverable              bool              `json:"discoverable" bson:"discoverable"`
	Indexable                 bool              `json:"indexable" bson:"indexable"`
	AttributionDomains        []string          `json:"attributionDomains,omitempty" bson:"attributionDomains"` // articles.hereus.net, and user URLs
	Attachment                []ActorAttachment `json:"attachment,omitempty" bson:"attachment,omitempty"`
	Image                     ActorIcon         `json:"image,omitempty" bson:"image,omitempty"`
}

type ActivityPubAttachment struct {
	Type      string `json:"type" bson:"type"` // "Image", "Video", "Audio", "Document", etc.
	MediaType string `json:"mediaType" bson:"mediaType"`
	Url       string `json:"url" bson:"url"`
	Name      string `json:"name,omitempty" bson:"name"`
}

type ActivityPubTag struct {
	Type string `json:"type" bson:"type"` // "Mention", "Hashtag", etc.
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Href string `json:"href" bson:"href"`
}

type ActivityPubNote struct {
	LdContext    interface{}             `json:"@context,omitempty" bson:"@context,omitempty"`
	Id           string                  `json:"id" bson:"id"`
	Type         string                  `json:"type" bson:"type,omitempty"`
	Summary      string                  `json:"summary,omitempty" bson:"summary,omitempty"`
	InReplyTo    string                  `json:"inReplyTo,omitempty" bson:"inReplyTo,omitempty"`
	Published    string                  `json:"published,omitempty" bson:"published,omitempty"`
	Updated      string                  `json:"updated,omitempty" bson:"updated,omitempty"`
	Url          string                  `json:"url,omitempty" bson:"url,omitempty"`
	CustomUrl    string                  `json:"hereus:url,omitempty" bson:"customUrl,omitempty"`
	AttributedTo string                  `json:"attributedTo,omitempty" bson:"attributedTo,omitempty"`
	To           interface{}             `json:"to,omitempty" bson:"to,omitempty"`
	Cc           interface{}             `json:"cc,omitempty" bson:"cc,omitempty"`
	Sensitive    bool                    `json:"sensitive,omitempty" bson:"sensitive,omitempty"` // content warning if true
	Content      string                  `json:"content" bson:"content,omitempty"`
	Attachment   []ActivityPubAttachment `json:"attachment,omitempty" bson:"attachment,omitempty"`
	Tag          []ActivityPubTag        `json:"tag,omitempty" bson:"tag,omitempty"` // mentions, hashtags, etc.
	Replies      interface{}             `json:"replies,omitempty" bson:"replies,omitempty"`
	Properties   map[string]interface{}  `json:"as:properties,omitempty" bson:"properties,omitempty"`
	PackageName  string                  `json:"as:packageName,omitempty" bson:"packageName,omitempty"`
}

type FollowModal struct {
	Who  string `bson:"who"`  // actor who follows
	Whom string `bson:"whom"` // actor who is followed
	Date string `bson:"date"` // YYYY-MM-DD HH:MM:SS
	Id   string `bson:"id"`
}

type ActivityPubAnnouncement struct {
	Id   string   `bson:"id"`
	To   string   `bson:"to"`
	Cc   []string `bson:"cc"`
	Who  string   `bson:"who"`
	What string   `bson:"what"`
	When string   `bson:"when"`
}
