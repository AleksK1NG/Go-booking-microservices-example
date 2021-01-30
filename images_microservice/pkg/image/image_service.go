package image

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"sync"

	"github.com/disintegration/gift"
)

// ImgToolkit Image toolkit service
type ImgToolkit struct {
	pool *sync.Pool
	g    *gift.GIFT
	buf  *bytes.Buffer
}

func (i *ImgToolkit) get() *gift.GIFT {
	g, ok := i.pool.Get().(*gift.GIFT)
	if !ok {
		return nil
	}
	return g
}

func (i *ImgToolkit) put(g *gift.GIFT) {
	i.pool.Put(g)
}

func NewImgToolkit(width, height int, resampling gift.Resampling) *ImgToolkit {
	pool := &sync.Pool{New: func() interface{} {
		g := gift.New(
			// gift.Resize(1024, 0, gift.LanczosResampling),
			gift.Resize(width, height, resampling),
			gift.Contrast(20),
			gift.Brightness(7),
			gift.Gamma(0.5),
			// gift.CropToSize(1024, 1024, gift.CenterAnchor),
		)

		return g
	}}

	return &ImgToolkit{pool: pool, buf: &bytes.Buffer{}}
}

func (i *ImgToolkit) ResizeImage(file io.Reader) ([]byte, error) {
	src, fileType, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	g := i.get()

	dst := image.NewNRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	defer i.buf.Reset()
	switch fileType {
	case "png":
		err = png.Encode(i.buf, dst)
		if err != nil {
			return nil, err
		}
	case "jpeg":
		err = jpeg.Encode(i.buf, dst, nil)
		if err != nil {
			return nil, err
		}
	case "jpg":
		err = jpeg.Encode(i.buf, dst, nil)
		if err != nil {
			return nil, err
		}
	case "gif":
		err = gif.Encode(i.buf, dst, nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid image format")
	}

	return i.buf.Bytes(), nil
}
