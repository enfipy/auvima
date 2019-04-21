package config

/* Settings */

type Settings struct {
	Coub      Coub      `yaml:"coub"`
	Instagram Instagram `yaml:"instagram"`
	Storage   Storage   `yaml:"storage"`
	Video     Video     `yaml:"video"`
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

// Instagram
type (
	Instagram struct {
		CredsPath            string   `yaml:"creds_path"`
		Creds                Creds    `yaml:"creds"`
		SuitabilityHours     uint64   `yaml:"suitability_hours"`
		MaterialAccounts     []string `yaml:"material_accounts"`
		MaterialCountToFetch uint32   `yaml:"material_count_to_fetch"`
	}

	Creds struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
)

// Storage
type Storage struct {
	Temporary  string `yaml:"tmp"`
	Finished   string `yaml:"finished"`
	Static     string `yaml:"static"`
	Production string `yaml:"prod"`
}

// Video
type (
	Video struct {
		Length  int64   `yaml:"length"`
		Timings Timings `yaml:"timings"`
	}

	Timings struct {
		FetchMaterial           string `yaml:"fetch_material"`
		GenerateProductionVideo string `yaml:"generate_production_video"`
		UploadVideo             string `yaml:"upload_video"`
	}
)
