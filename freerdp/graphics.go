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

	return C.TRUE
}

//export webRdpBitmapSetSurface
func webRdpBitmapSetSurface(context *C.rdpContext, bitmap *C.rdpBitmap, primary C.BOOL) C.BOOL {

	return C.TRUE
}
