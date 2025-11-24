// internal/data/prices.go
package data

import (
	"log"
	"os"
)

var PriceMap map[string]string

func InitPriceMap() {
	PriceMap = map[string]string{
		"Banner":         os.Getenv("BannerID"),
		"Thumbnail":      os.Getenv("YtThumbnailID"),
		"PFP":            os.Getenv("ProfilePictureID"),
		"Twitch Emotes":  os.Getenv("TwitchEmotesID"),
		"Logos":          os.Getenv("LogosID"),
		"Streamer Packs": os.Getenv("StreamerPacksID"),
	}

	log.Println("PriceMap loaded:", PriceMap)
}
