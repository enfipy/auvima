package models

import "time"

type (
	Link struct {
		URL  string `json:"url"`
		Size int32  `json:"size"`
	}

	Media struct {
		High *Link `json:"high"`
		Med  *Link `json:"med"`
	}

	HTML5 struct {
		Video Media `json:"video"`
		Audio Media `json:"audio"`
	}

	FileVersions struct {
		HTML5 *HTML5 `json:"html5"`
	}

	Channel struct {
		Id             int    `json:"id"`
		Permalink      string `json:"permalink"`
		Description    string `json:"description"`
		Title          string `json:"title"`
		FollowersCount int    `json:"followers_count"`
		FollowingCount int    `json:"following_count"`
	}

	Tag struct {
		Id    int    `json:"id"`
		Value string `json:"value"`
	}

	Coub struct {
		Id            int           `json:"id"`
		Type          string        `json:"type"`
		Permalink     string        `json:"permalink"`
		Title         string        `json:"title"`
		ChannelId     int           `json:"channel_id"`
		CreatedAt     *time.Time    `json:"created_at"`
		UpdatedAt     *time.Time    `json:"updated_at"`
		Duration      float64       `json:"duration"`
		ViewsCount    int           `json:"views_count"`
		OriginalSound bool          `json:"original_sound"`
		HasSound      bool          `json:"has_sound"`
		FileVersions  *FileVersions `json:"file_versions"`
		AgeRestricted bool          `json:"age_restricted"`
		AllowReuse    bool          `json:"allow_reuse"`
		Banned        bool          `json:"banned"`
		Channel       *Channel      `json:"channel"`
		Tags          []Tag         `json:"tags"`
	}
)
