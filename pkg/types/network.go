package types

type SubscriptionPlan struct {
	Name    string `json:"name" bson:"name"`
	Storage int    `json:"storage" bson:"storage"`
}

type NetworkOS struct {
	Arch    string `json:"arch"`
	Family  string `json:"family"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type NetworkRules struct {
	NewAccountsAllowed bool `json:"newAccountsAllowed"`
}

type NetworkSoftware struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	Version string `json:"version"`
	Build   int    `json:"build"`
}

type Network struct {
	TheProtocolsVersion float64            `json:"version"`
	HelpUsername        string             `json:"help"`
	SubscriptionPlans   []SubscriptionPlan `json:"subscriptionPlans"`
	OS                  NetworkOS          `json:"os"`
	Rules               NetworkRules       `json:"rules"`
	Software            NetworkSoftware    `json:"software"`
	AuthorizationUrl    string             `json:"authorizationUrl"`
}
