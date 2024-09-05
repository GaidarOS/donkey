package thumbnails

import (
	"donkey/config"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/nfnt/resize"
)

func GenerateThumbnailFromImage(filepath string, savefoldername string) error {
	image_file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer image_file.Close()

	filetype := strings.Split(path.Base(filepath), ".")[1]
	var my_image image.Image
	switch filetype {
	case "jpeg":
		my_image, err = jpeg.Decode(image_file)
	case "jpg":
		my_image, err = jpeg.Decode(image_file)
	case "png":
		my_image, err = png.Decode(image_file)
	}

	if err != nil {
		return err
	}

	newImage := resize.Resize(150, 150, my_image, resize.Lanczos3)

	output_file, outputErr := os.Create(path.Join(config.AppConf.Dir, savefoldername, path.Base(filepath)))
	if outputErr != nil {
		return outputErr
	}
	jpeg.Encode(output_file, newImage, &jpeg.Options{Quality: jpeg.DefaultQuality})
	defer output_file.Close()
	return nil
}

func GenerateThumbnailFromPdf(filepath string, savefoldername string) error {
	doc, err := fitz.New(filepath)
	if err != nil {
		return err
	}

	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(config.AppConf.Dir, savefoldername, strings.Replace(path.Base(filepath), ".pdf", ".jpg", 1)))
	if err != nil {
		return err
	}

	newImage := resize.Resize(150, 150, img, resize.Lanczos3)

	err = jpeg.Encode(f, newImage, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}
