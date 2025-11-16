package gobs

//#include <obs/obs.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type OutputFormat uint32

const (
	VIDEO_FORMAT_NONE OutputFormat = iota

	/* planar 4:2:0 formats */
	VIDEO_FORMAT_I420 /* three-plane */
	VIDEO_FORMAT_NV12 /* two-plane, luma and packed chroma */

	/* packed 4:2:2 formats */
	VIDEO_FORMAT_YVYU
	VIDEO_FORMAT_YUY2 /* YUYV */
	VIDEO_FORMAT_UYVY

	/* packed uncompressed formats */
	VIDEO_FORMAT_RGBA
	VIDEO_FORMAT_BGRA
	VIDEO_FORMAT_BGRX
	VIDEO_FORMAT_Y800 /* grayscale */

	/* planar 4:4:4 */
	VIDEO_FORMAT_I444

	/* more packed uncompressed formats */
	VIDEO_FORMAT_BGR3

	/* planar 4:2:2 */
	VIDEO_FORMAT_I422

	/* planar 4:2:0 with alpha */
	VIDEO_FORMAT_I40A

	/* planar 4:2:2 with alpha */
	VIDEO_FORMAT_I42A

	/* planar 4:4:4 with alpha */
	VIDEO_FORMAT_YUVA

	/* packed 4:4:4 with alpha */
	VIDEO_FORMAT_AYUV

	/* planar 4:2:0 format, 10 bpp */
	VIDEO_FORMAT_I010 /* three-plane */
	VIDEO_FORMAT_P010 /* two-plane, luma and packed chroma */

	/* planar 4:2:2 format, 10 bpp */
	VIDEO_FORMAT_I210

	/* planar 4:4:4 format, 12 bpp */
	VIDEO_FORMAT_I412

	/* planar 4:4:4:4 format, 12 bpp */
	VIDEO_FORMAT_YA2L

	/* planar 4:2:2 format, 16 bpp */
	VIDEO_FORMAT_P216 /* two-plane, luma and packed chroma */

	/* planar 4:4:4 format, 16 bpp */
	VIDEO_FORMAT_P416 /* two-plane, luma and packed chroma */

	/* packed 4:2:2 format, 10 bpp */
	VIDEO_FORMAT_V210

	/* packed uncompressed 10-bit format */
	VIDEO_FORMAT_R10L
)

type VideoColorSpace uint32

const (
	VIDEO_CS_DEFAULT VideoColorSpace = iota
	VIDEO_CS_601
	VIDEO_CS_709
	VIDEO_CS_SRGB
	VIDEO_CS_2100_PQ
	VIDEO_CS_2100_HLG
)

type VideoRangeType uint32

const (
	VIDEO_RANGE_DEFAULT VideoRangeType = iota
	VIDEO_RANGE_PARTIAL
	VIDEO_RANGE_FULL
)

type ScaleType uint32

const (
	VIDEO_SCALE_DEFAULT ScaleType = iota
	VIDEO_SCALE_POINT
	VIDEO_SCALE_FAST_BILINEAR
	VIDEO_SCALE_BILINEAR
	VIDEO_SCALE_BICUBIC
)

type GraphicsModule *uint8

var (
	GraphicsModuleOpenGL GraphicsModule = unsafe.StringData("libobs-opengl\x00")
	GraphicsModuleD3D11  GraphicsModule = unsafe.StringData("libobs-d3d11\x00")
)

type VideoInfo struct {
	GraphicsModule GraphicsModule
	FPSNum         uint32
	FPSDen         uint32
	BaseWidth      uint32
	BaseHeight     uint32
	OutputWidth    uint32
	OutputHeight   uint32
	OutputFormat   OutputFormat
	Adapter        uint32
	GPUConversion  bool
	ColorSpace     VideoColorSpace
	Range          VideoRangeType
	ScaleType      ScaleType
}

func ResetVideo(info *VideoInfo) error {
	// #cgo noescape obs_reset_video
	// #cgo nocallback obs_reset_video
	err := C.obs_reset_video((*C.struct_obs_video_info)(unsafe.Pointer(info)))
	switch err {
	case 0:
		return nil // OBS_VIDEO_SUCCESS
	case -1:
		return errors.New("video failed")
	case -2:
		return errors.New("video not supported")
	case -3:
		return errors.New("invalid parameter")
	case -4:
		return errors.New("video already active")
	case -5:
		return errors.New("video module not found")
	default:
		return fmt.Errorf("Unknown error: %v", err)
	}
}

type Video struct {
	c *C.video_t
}
