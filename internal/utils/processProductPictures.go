package utils

import (
	"os"
	"shop-dashboard/internal/models"
)

func ProcessPictures(urls []string, name string) []models.ProductPicture {
	var pictures []models.ProductPicture
	if urls != nil {
		for _, imgURL := range urls {
			pictures = append(pictures, models.ProductPicture{
				ImageUrl:      imgURL,
				Title:         name,
				AlternateText: name,
			})
		}

		return pictures
	}

	return []models.ProductPicture{{
		ImageUrl:      os.Getenv("CLIENT_SERVER") + "/images/no_image_placeholder.jpg",
		Title:         "No Image",
		AlternateText: "No Image",
	}}
}
