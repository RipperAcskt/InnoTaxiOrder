package model

type ElasticModel struct {
	Hits struct {
		Hits []struct {
			ID     string `json:"_id"`
			Source Order  `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

const (
	Econom   string = "econom"
	Comfort  string = "comfort"
	Business string = "business"

	User   string = "user"
	Driver string = "driver"
)
