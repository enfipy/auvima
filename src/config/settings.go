package config

/* Settings */

type Settings struct {
	Coub    Coub    `yaml:"coub"`
	Storage Storage `yaml:"storage"`
}

/* Inner structs */

// Coub
type (
	Coub struct {
		Secrets CoubSecrets `yaml:"secrets"`
		Urls    CoubUrls    `yaml:"urls"`
	}

	CoubSecrets struct {
		AccessToken string `yaml:"access_token"`
	}

	CoubUrls struct {
		BaseUrl string `yaml:"base_url"`
	}
)

// Coub
type Storage struct {
	Temporary string `yaml:"tmp"`
	Finished  string `yaml:"finished"`
	Static    string `yaml:"static"`
}
