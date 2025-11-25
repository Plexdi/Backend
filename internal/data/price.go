// internal/data/prices.go
package data

import (
	"log"
	"os"
)

var PriceMap map[string][]string

func InitPriceMap() {
	PriceMap := map[string][]string{
		"Banner": []string{
			os.Getenv("StarterBannerPriceID"),
			os.Getenv("StandardBannerPriceID"),
			os.Getenv("PremiumBannerPriceID"),
		},
		"Thumbnail": []string{
			os.Getenv("StarterThumbnailPriceID"),
			os.Getenv("StandardThumbnailPriceID"),
			os.Getenv("PremiumThumbnailPriceID"),
		},
		"PFP": []string{
			os.Getenv("StarterPFPPriceID"),
			os.Getenv("StandardPFPPriceID"),
			os.Getenv("PremiumPFPPriceID"),
		},
		"Twitch Emotes": []string{
			os.Getenv("StarterTwitchEmotesPriceID"),
			os.Getenv("StandardTwitchEmotesPriceID"),
			os.Getenv("PremiumTwitchEmotesPriceID"),
		},
		"Logos": []string{
			os.Getenv("StarterLogosPriceID"),
			os.Getenv("StandardLogosPriceID"),
			os.Getenv("PremiumLogosPriceID"),
		},
		"Streamer Packs": []string{
			os.Getenv("StarterStreamerPacksPriceID"),
			os.Getenv("StandardStreamerPacksPriceID"),
			os.Getenv("PremiumStreamerPacksPriceID"),
		},
	}

	log.Println("PriceMap loaded:", PriceMap)
}
