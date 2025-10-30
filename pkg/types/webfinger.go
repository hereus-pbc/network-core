package types

type WebFingerLinks struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

type WebFingerResponse struct {
	Subject string           `json:"subject" bson:"subject"`
	Aliases []string         `json:"aliases" bson:"aliases"`
	Links   []WebFingerLinks `json:"links" bson:"links"`
}
