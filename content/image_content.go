package content

import (
	"bytes"
	"fmt"
	"image"
	"sync"
)

type Image struct {
	Image image.Image

	SourceFormat string
}

type imageContent struct {
	*content
	imageCache Image
	once       sync.Once
}

func (i *imageContent) Image() (Image, error) {
	var err error
	i.once.Do(func() {
		var r *bytes.Reader
		r, err = i.BytesReader()
		if err != nil {
			return
		}
		im, format, err := image.Decode(r)
		if err != nil {
			return
		}
		i.imageCache = Image{
			Image:        im,
			SourceFormat: format,
		}
	})
	if err != nil {
		return Image{}, err
	}
	return i.imageCache, nil
}

func (i *imageContent) BytesReader() (*bytes.Reader, error) {
	if i.content != nil {
		return i.content.BytesReader()
	} else {
		return nil, fmt.Errorf("ImageContent has no underlying buffer")
	}
}
