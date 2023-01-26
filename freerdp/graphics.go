package freerdp

/*
#if __APPLE__
#include <TargetConditionals.h>
#endif

#include "freerdp/freerdp.h"

extern BOOL webRdpBitmapNew(rdpContext* context, rdpBitmap* bitmap);
extern void webRdpBitmapFree(rdpContext* context, rdpBitmap* bitmap);
extern BOOL webRdpBitmapPaint(rdpContext* context, rdpBitmap* bitmap);
extern BOOL webRdpBitmapSetSurface(rdpContext* context, rdpBitmap* bitmap, BOOL primary);

static void init_bitmap(rdpGraphics *graphics)
{
	rdpBitmap bitmap = *graphics->Bitmap_Prototype;
  	bitmap.size = sizeof(rdpBitmap);
    bitmap.New = webRdpBitmapNew;
    bitmap.Free = webRdpBitmapFree;
    bitmap.Paint = webRdpBitmapPaint;
    bitmap.SetSurface = webRdpBitmapSetSurface;
    graphics_register_bitmap(graphics, &bitmap);
}
*/
import "C"
import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"unsafe"
)

func registerGraphics(graphics *C.rdpGraphics) {
	C.init_bitmap(graphics)
}

//export webRdpBitmapNew
func webRdpBitmapNew(context *C.rdpContext, bitmap *C.rdpBitmap) C.BOOL {
	return C.TRUE
}

//export webRdpBitmapFree
func webRdpBitmapFree(context *C.rdpContext, bitmap *C.rdpBitmap) {
}

//export webRdpBitmapPaint
func webRdpBitmapPaint(context *C.rdpContext, bitmap *C.rdpBitmap) C.BOOL {
	if bitmap.data == nil {
		return C.TRUE
	}

	////fmt.Printf("left: %v, top: %v, right: %v, bottom: %v, width: %v, height: %v, format: %v, flags: %v, length: %v, comparessed: %v, ep")
	//fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
	//	bitmap.left, bitmap.top, bitmap.right, bitmap.bottom, bitmap.width, bitmap.height,
	//	bitmap.format, bitmap.flags, bitmap.length, bitmap.compressed, bitmap.ephemeral)

	//bitmap.format = 0x20010888
	//bitmap.flags = 0

	w := int(bitmap.width)
	h := int(bitmap.height)
	data := C.GoBytes(unsafe.Pointer(bitmap.data), C.int(bitmap.length))

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y, i := 0, 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: data[i+1], G: data[i+2], B: data[i+3], A: 255})
			i += 4
		}
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return C.TRUE
	}

	Try(func() { // Prevent `panic: send on closed channel`
		client := getClientFromContent(context)
		client.dataChan <- &Message{Bitmap: &Bitmap{
			X:    int(bitmap.left),
			Y:    int(bitmap.top),
			W:    w,
			H:    h,
			Data: buf.Bytes(),
		}}
	})

	return C.TRUE
}

//export webRdpBitmapSetSurface
func webRdpBitmapSetSurface(context *C.rdpContext, bitmap *C.rdpBitmap, primary C.BOOL) C.BOOL {
	return C.TRUE
}
