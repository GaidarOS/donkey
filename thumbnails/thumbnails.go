package thumbnails

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"receipt_store/config"
	"strings"

	"github.com/gen2brain/go-fitz"
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

	my_sub_image := my_image.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 0, 150, 150))

	output_file, outputErr := os.Create("output.jpeg")
	if outputErr != nil {
		return outputErr
	}
	jpeg.Encode(output_file, my_sub_image, &jpeg.Options{Quality: jpeg.DefaultQuality})
	return nil
}

func GenerateThumbnailFromPdf(filepath string, savefoldername string) error {
	doc, err := fitz.New(filepath)
	if err != nil {
		return err
	}

	defer doc.Close()

	img,err := doc.Image(0)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(config.AppConf.Dir, savefoldername, strings.Replace(path.Base(filepath), ".pdf", ".png", 1)))
	if err != nil {
		return err
	}

	my_sub_image := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 0, 150, 150))

	err = jpeg.Encode(f, my_sub_image, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		return err
	}

	f.Close()
	return nil
}
