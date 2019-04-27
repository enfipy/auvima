package config

/* Settings */

type Settings struct {
	Coub      Coub      `yaml:"coub"`
	Instagram Instagram `yaml:"instagram"`
	Youtube   Youtube   `yaml:"youtube"`
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
		CredsPath            string         `yaml:"creds_path"`
		Creds                InstagramCreds `yaml:"creds"`
		SuitabilityHours     uint64         `yaml:"suitability_hours"`
		MaterialAccounts     []string       `yaml:"material_accounts"`
		MaterialCountToFetch uint32         `yaml:"material_count_to_fetch"`
	}

	InstagramCreds struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
)

// Youtube
type (
	Youtube struct {
		Creds YoutubeCreds `yaml:"creds"`
	}

	YoutubeCreds struct {
		AccessToken  string `yaml:"access_token"`
		TokenType    string `yaml:"token_type"`
		RefreshToken string `yaml:"refresh_token"`
		Expiry       uint64 `yaml:"expiry"`
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
		Title       string  `yaml:"title"`
		CategoryId  string  `yaml:"category_id"`
		Privacy     string  `yaml:"privacy"`
		Tags        string  `yaml:"tags"`
		Description string  `yaml:"description"`
		Length      int64   `yaml:"length"`
		FrameLength int64   `yaml:"frame_length"`
		Timings     Timings `yaml:"timings"`
	}

	Timings struct {
		FetchMaterial           string `yaml:"fetch_material"`
		GenerateProductionVideo string `yaml:"generate_production_video"`
		UploadVideo             string `yaml:"upload_video"`
	}
)
