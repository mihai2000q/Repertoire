package model

type ClientCredentials struct {
	GrantType    string
	ClientID     string
	ClientSecret string
}

func NewClientCredentials(grantType string, clientID string, clientSecret string) ClientCredentials {
	return ClientCredentials{
		GrantType:    grantType,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}
