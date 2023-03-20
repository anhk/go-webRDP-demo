package freerdp

/*
#if __APPLE__
#include <TargetConditionals.h>
#endif

#include "freerdp/freerdp.h"

typedef const BYTE cByte;
extern BOOL webRdpBitmapNew(rdpContext* context, rdpBitmap* bitmap);
extern void webRdpBitmapFree(rdpContext* context, rdpBitmap* bitmap);
extern BOOL webRdpBitmapPaint(rdpContext* context, rdpBitmap* bitmap);
extern BOOL webRdpBitmapSetSurface(rdpContext* context, rdpBitmap* bitmap, BOOL primary);
extern BOOL webRdpBitmapDecompress(rdpContext* context, rdpBitmap* bitmap, cByte* data,
	                                   UINT32 width, UINT32 height, UINT32 bpp, UINT32 length,
	                                   BOOL compressed, UINT32 codec_id);

static void init_bitmap(rdpGraphics *graphics)
{
	rdpBitmap bitmap = *graphics->Bitmap_Prototype;
  	bitmap.size = sizeof(rdpBitmap);
    bitmap.New = webRdpBitmapNew;
    bitmap.Free = webRdpBitmapFree;
    bitmap.Paint = webRdpBitmapPaint;
    bitmap.SetSurface = webRdpBitmapSetSurface;
	bitmap.Decompress = webRdpBitmapDecompress;
    graphics_register_bitmap(graphics, &bitmap);
}

typedef const rdpGlyph cRdpGlyph;
extern BOOL webRdpGlyphNew(rdpContext* context, cRdpGlyph* glyph);
extern void webRdpGlyphFree(rdpContext* context, rdpGlyph* glyph);
extern BOOL webRdpGlyphDraw(rdpContext* context, cRdpGlyph* glyph, INT32 x, INT32 y,
    INT32 w, INT32 h, INT32 sx, INT32 sy, BOOL fOpRedundant);
extern BOOL webRdpGlyphBeginDraw(rdpContext* context, INT32 x, INT32 y, INT32 width,
	INT32 height, UINT32 bgcolor, UINT32 fgcolor, BOOL fOpRedundant);
extern BOOL webRdpGlyphEndDraw(rdpContext* context, INT32 x, INT32 y, INT32 width, INT32 height,
    UINT32 bgcolor, UINT32 fgcolor);

static void init_glyph(rdpGraphics *graphics)
{
	rdpGlyph glyph = *graphics->Glyph_Prototype;
    glyph.size = sizeof(rdpGlyph);
    glyph.New = webRdpGlyphNew;
    glyph.Free = webRdpGlyphFree;
    glyph.Draw = webRdpGlyphDraw;
    glyph.BeginDraw = webRdpGlyphBeginDraw;
    glyph.EndDraw = webRdpGlyphEndDraw;
    graphics_register_glyph(graphics, &glyph);
}
typedef const rdpPointer cRdpPointer;
extern BOOL webRdpPointerNew(rdpContext* context, rdpPointer* pointer);
extern void webRdpPointerFree(rdpContext* context, rdpPointer* pointer);
extern BOOL webRdpPointerSet(rdpContext* context, cRdpPointer* pointer);
extern BOOL webRdpPointerSetNull(rdpContext* context);
extern BOOL webRdpPointerSetDefault(rdpContext* context);

static void init_pointer(rdpGraphics *graphics)
{
    rdpPointer pointer = *graphics->Pointer_Prototype;
    pointer.size = sizeof(rdpPointer);
    pointer.New = webRdpPointerNew;
    pointer.Free = webRdpPointerFree;
    pointer.Set = webRdpPointerSet;
    pointer.SetNull = webRdpPointerSetNull;
    pointer.SetDefault = webRdpPointerSetDefault;
    graphics_register_pointer(graphics, &pointer);
}
*/
import "C"
import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"runtime"
	"unsafe"

	"github.com/GoFeGroup/go-bitmap/glog"

	gobitmap "github.com/GoFeGroup/go-bitmap"
)

func registerGraphics(graphics *C.rdpGraphics) {
	C.init_bitmap(graphics)
	C.init_glyph(graphics)
	C.init_pointer(graphics)
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

//export webRdpBitmapDecompress
func webRdpBitmapDecompress(context *C.rdpContext, bitmap *C.rdpBitmap, data *C.cByte,
	width, height, bpp, length C.UINT32, compressed C.BOOL, codec_id C.UINT32) C.BOOL {
	//fmt.Printf("%+v\n", bitmap)
	//fmt.Println("w,h,bpp,length:", width, height, bpp, length)

	d := C.GoBytes(unsafe.Pointer(data), C.int(length))

	//fmt.Println("compressed:", compressed, "codecId:", codec_id)

	var bmp *gobitmap.BitMap
	option := &gobitmap.Option{
		Width:       int(width),
		Height:      int(height),
		BitPerPixel: int(bpp),
		Data:        d,
	}

	glog.SetLevel(glog.NONE)
	if compressed != C.TRUE {
		bmp = gobitmap.NewBitmapFromPlain(option)
	} else if bpp != 32 {
		bmp = gobitmap.NewBitmapFromRLE(option)
	} else {
		bmp = gobitmap.NewBitMapFromRDP6(option)
	}
	//fmt.Println(bmp.ToPng())

	Try(func() { // Prevent `panic: send on closed channel`
		client := getClientFromContent(context)
		client.dataChan <- &Message{Bitmap: &Bitmap{
			X:    int(bitmap.left),
			Y:    int(bitmap.top),
			W:    int(width),
			H:    int(height),
			Data: bmp.ToPng(),
		}}
	})

	return C.TRUE
}

//export webRdpGlyphNew
func webRdpGlyphNew(context *C.rdpContext, glyph *C.cRdpGlyph) C.BOOL {
	fmt.Println(runtime.Caller(0))
	return C.TRUE
}

//export webRdpGlyphFree
func webRdpGlyphFree(context *C.rdpContext, glyph *C.rdpGlyph) {
	fmt.Println(runtime.Caller(0))
}

//export webRdpGlyphDraw
func webRdpGlyphDraw(context *C.rdpContext, glyph *C.cRdpGlyph, x, y,
	w, h, sx, sy C.INT32, fOpRedundant C.BOOL) C.BOOL {
	fmt.Println(runtime.Caller(0))
	return C.TRUE
}

//export webRdpGlyphBeginDraw
func webRdpGlyphBeginDraw(context *C.rdpContext, x, y, width,
	height C.INT32, bgcolor, fgcolor C.UINT32, fOpRedundant C.BOOL) C.BOOL {
	fmt.Println(runtime.Caller(0))
	return C.TRUE
}

//export webRdpGlyphEndDraw
func webRdpGlyphEndDraw(context *C.rdpContext, x, y, width, height C.INT32,
	bgcolor, fgcolor C.UINT32) C.BOOL {
	fmt.Println(runtime.Caller(0))
	return C.TRUE
}

//export webRdpPointerNew
func webRdpPointerNew(context *C.rdpContext, pointer *C.rdpPointer) C.BOOL {
	//fmt.Println(runtime.Caller(0))
	return C.TRUE
}

//export webRdpPointerFree
func webRdpPointerFree(context *C.rdpContext, pointer *C.rdpPointer) {
	//fmt.Println(runtime.Caller(0))
}

//export webRdpPointerSet
func webRdpPointerSet(context *C.rdpContext, pointer *C.cRdpPointer) C.BOOL {
	//fmt.Println(runtime.Caller(0))
	return C.TRUE
}

//export webRdpPointerSetNull
func webRdpPointerSetNull(context *C.rdpContext) C.BOOL {
	fmt.Println(runtime.Caller(0))
	return C.TRUE
}

//export webRdpPointerSetDefault
func webRdpPointerSetDefault(context *C.rdpContext) C.BOOL {
	fmt.Println(runtime.Caller(0))
	return C.TRUE
}
