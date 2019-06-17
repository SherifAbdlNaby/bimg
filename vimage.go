package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

import (
	"runtime"
)

type VipsImage struct {
	image     *C.VipsImage
	srcBuf    []byte
	imageType ImageType
}

func NewVipsImage(buf []byte) (*VipsImage, error) {

	var err error

	ref := &VipsImage{}

	ref.srcBuf = buf

	ref.image, ref.imageType, err = loadImage(ref.srcBuf)

	runtime.SetFinalizer(ref, func(ref *VipsImage) {
		defer C.g_object_unref(C.gpointer(ref.image))
	})

	return ref, err
}

func (i *VipsImage) Process(o Options) error {
	//defer C.vips_thread_shutdown()
	o = applyDefaults(o, i.imageType)

	image, err := process(i.image, i.imageType, o, &i.srcBuf)
	if err != nil {
		return err
	}

	i.image = image

	return nil
}

func (i *VipsImage) Save(o Options) ([]byte, error) {
	o = applyDefaults(o, i.imageType)

	saveOptions := vipsSaveOptions{
		Quality:        o.Quality,
		Type:           o.Type,
		Compression:    o.Compression,
		Interlace:      o.Interlace,
		NoProfile:      o.NoProfile,
		Interpretation: o.Interpretation,
		OutputICC:      o.OutputICC,
		StripMetadata:  o.StripMetadata,
		Lossless:       o.Lossless,
	}
	// Finally get the resultant buffer
	buf, image, err := vipsSave(i.image, saveOptions)
	if err != nil {
		return nil, err
	}

	i.image = image

	return buf, nil
}

func (i *VipsImage) Clone() *VipsImage {

	clone := &VipsImage{
		image:     vipsCopy(i.image),
		srcBuf:    i.srcBuf,
		imageType: i.imageType,
	}

	runtime.SetFinalizer(clone, func(ref *VipsImage) {
		defer C.g_object_unref(C.gpointer(ref.image))
	})

	return clone
}
