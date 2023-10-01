package content

import (
	"bytes"
	"fmt"
	"image"
	"sync"
)

type imageContent struct {
	*content
	imageCache image.Image
	once       sync.Once
}

func (i *imageContent) Image() (im image.Image, err error) {
	im = i.imageCache
	i.once.Do(func() {
		var r *bytes.Reader
		r, err = i.BytesReader()
		if err != nil {
			return
		}
		im, _, err = image.Decode(r)
		if err != nil {
			return
		}
		i.imageCache = im
	})
	return
}

func (i *imageContent) BytesReader() (*bytes.Reader, error) {
	if i.content != nil {
		return i.content.BytesReader()
	} else {
		return nil, fmt.Errorf("ImageContent has no underlying buffer")
	}
}
