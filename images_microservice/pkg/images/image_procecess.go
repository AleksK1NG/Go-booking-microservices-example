package images

import (
	"bytes"

	"github.com/disintegration/gift"
)

type ImgResizer struct {
	Gift   *gift.GIFT
	Buffer *bytes.Buffer
}

func NewImgResizer(filters ...gift.Filter) *ImgResizer {
	g := gift.New(
		filters...,
	// gift.Resize(1024, 0, gift.LanczosResampling),
	// gift.Resize(width, height, resampling),
	// gift.Contrast(20),
	// gift.Brightness(7),
	// gift.Gamma(0.5),
	// gift.CropToSize(1024, 1024, gift.CenterAnchor),
	)
	b := &bytes.Buffer{}
	return &ImgResizer{Gift: g, Buffer: b}
}
