package api

import "github.com/cloudinary/cloudinary-go/v2"

var cld *cloudinary.Cloudinary

func LoadCloudinary(cloudName, apiKey, apiSecret string) error {
	var err error

	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return err
	}

	return nil
}
