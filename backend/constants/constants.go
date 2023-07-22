package constants

import (
	"time"
)

const (
	// TTL otp redis 5 min
	TIME_TO_LIVE = 5 * time.Minute

	// Google Cloud Storage base URL
	STORAGE_URL = "https://storage.googleapis.com"

	// pdf or assignments dir
	ASSIGNMENTS_DIR = "assignments"

	// videos dir
	VIDEOS_DIR = "videos"

	// images dir
	IMAGES_DIR = "images"
)
