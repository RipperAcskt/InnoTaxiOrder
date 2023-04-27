package model

type (
	UserType  string
	ClassType string
)

type ElasticModel struct {
	Hits struct {
		Hits []struct {
			ID     string `json:"_id"`
			Source Order  `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

const (
	Econom   ClassType = "econom"
	Comfort  ClassType = "comfort"
	Business ClassType = "business"

	User   UserType = "user"
	Driver UserType = "driver"
)

func NewUserType(s string) UserType {
	return UserType(s)
}

func NewClassType(s string) ClassType {
	return ClassType(s)
}
