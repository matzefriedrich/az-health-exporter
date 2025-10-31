package resources

import (
	"embed"
)

//go:embed banner.txt
var Resources embed.FS

const (
	BannerTxt = "banner.txt"
)
