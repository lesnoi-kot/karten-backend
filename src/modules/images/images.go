package images

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"

	_ "image/jpeg"

	_ "golang.org/x/image/webp"

	"github.com/lesnoi-kot/karten-backend/src/settings"
	"github.com/nfnt/resize"
)

type ImageInfo struct {
	image.Config
	MIMEType string
}

func ParseImage(data io.Reader) (*ImageInfo, error) {
	imageConfig, imageType, err := image.DecodeConfig(data)
	if err != nil {
		return nil, fmt.Errorf("Image decoding error: %w", err)
	}

	switch imageType {
	case "png":
	case "jpeg":
	case "webp":
	default:
		return nil, fmt.Errorf("Unsupported image type: %s", imageType)
	}

	imageInfo := ImageInfo{
		Config:   imageConfig,
		MIMEType: "image/" + imageType,
	}

	return &imageInfo, nil
}

func MakeProjectAvatarThumbnail(data io.Reader) (io.Reader, error) {
	img, _, err := image.Decode(data)
	if err != nil {
		return nil, err
	}

	thumb := resize.Resize(
		settings.Projects.AvatarThumbnailSize,
		settings.Projects.AvatarThumbnailSize,
		img,
		resize.Bilinear,
	)
	buff := bytes.NewBuffer(make([]byte, 0, 50*1024))

	if err = png.Encode(buff, thumb); err != nil {
		return nil, err
	}

	return buff, nil
}
