package thumbnails

import (
	"donkey/config"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/nfnt/resize"
)

func GenerateThumbnailFromImage(filepath string, savefoldername string) error {
	image_file, err := os.Open(filepath)
	if err != nil {
		slog.Error("Error while opening the image file", err)
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
	default:
		// In the odd case an image comes through that is none of the above
		err = errors.New("image wasn't one of png/jpg/jpeg")
	}

	if err != nil {
		slog.Error("Error while decoding image", err)
		return err
	}

	newImage := resize.Resize(150, 150, my_image, resize.Lanczos3)

	output_file, outputErr := os.Create(path.Join(config.AppConf.Dir, savefoldername, path.Base(filepath)))
	if outputErr != nil {
		return outputErr
	}

	err = jpeg.Encode(output_file, newImage, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		slog.Error("Could not encode image", err)
		return err
	}
	defer output_file.Close()
	return nil
}

func GenerateThumbnailFromPdf(filepath string, savefoldername string) error {
	doc, err := fitz.New(filepath)
	if err != nil {
		slog.Error("Error while provisioning a thumbnail", err)
		return err
	}

	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		slog.Error("Error while generating a thumbnail", err)
		return err
	}

	f, err := os.Create(path.Join(config.AppConf.Dir, savefoldername, strings.Replace(path.Base(filepath), ".pdf", ".jpg", 1)))
	if err != nil {
		slog.Error("Error while creating a thumbnail", err)
		return err
	}

	newImage := resize.Resize(150, 150, img, resize.Lanczos3)

	err = jpeg.Encode(f, newImage, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		slog.Error("Error while encoding the thumbnail", err)
		return err
	}
	defer f.Close()

	return nil
}
