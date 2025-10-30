package types

type UserId struct {
	ProfilePhoto string `json:"profilePhoto"`
	FirstName    string `json:"firstName"`
	MiddleName   string `json:"middleName"`
	LastName     string `json:"lastName"`
	NameSuffix   string `json:"nameSuffix"`
	Birthdate    string `json:"birthdate"`
	Country      string `json:"country"`
	Gender       string `json:"gender"`
	PhoneNumber  string `json:"phoneNumber"`
	Timezone     string `json:"timezone"`
	Description  string `json:"description"`
	RsaPublicKey string `json:"rsaPublicKey"`
	SignatureId  string `json:"signatureId"`
	EncryptionId string `json:"encryptionId"`
}
