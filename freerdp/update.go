package freerdp

/*
#if __APPLE__
#include <TargetConditionals.h>
#endif

#include "freerdp/freerdp.h"

typedef const PLAY_SOUND_UPDATE cPlaySoundUpdate;
typedef const DSTBLT_ORDER cDstbltOrder;
typedef const OPAQUE_RECT_ORDER cOpaqueRectOrder;

extern BOOL webRdpUpdateBeginPaint(rdpContext* context);
extern BOOL webRdpUpdateEndPaint(rdpContext* context);
extern BOOL webRdpUpdatePlaySound(rdpContext* context, cPlaySoundUpdate* play_sound);
extern BOOL webRdpUpdateDesktopResize(rdpContext* context);
extern BOOL webRdpUpdateSetKeyboardIndicators(rdpContext* context, UINT16 led_flags);
extern BOOL webRdpUpdateSetKeyboardImeStatus(rdpContext* context, UINT16 imeId, UINT32 imeState, UINT32 imeConvMode);
extern BOOL webRdpDstBlt(rdpContext* context, cDstbltOrder* dstblt);
extern BOOL webRdpOpaqueRect(rdpContext* context, cOpaqueRectOrder* opaque_rect);

static void init_update_callbacks(rdpUpdate *update)
{
	rdpPrimaryUpdate *primary = update->primary;

    update->BeginPaint = webRdpUpdateBeginPaint;
    update->EndPaint = webRdpUpdateEndPaint;
    update->PlaySound = webRdpUpdatePlaySound;
    update->DesktopResize = webRdpUpdateDesktopResize;
    update->SetKeyboardIndicators = webRdpUpdateSetKeyboardIndicators;
    update->SetKeyboardImeStatus = webRdpUpdateSetKeyboardImeStatus;

	primary->DstBlt = webRdpDstBlt;
	primary->OpaqueRect = webRdpOpaqueRect;
}
*/
import "C"
import "fmt"

func registerUpdateCallbacks(update *C.rdpUpdate) {
	C.init_update_callbacks(update)
}

//export webRdpUpdateBeginPaint
func webRdpUpdateBeginPaint(context *C.rdpContext) C.BOOL {
	return C.TRUE
}

//export webRdpUpdateEndPaint
func webRdpUpdateEndPaint(context *C.rdpContext) C.BOOL {
	return C.TRUE
}

//export webRdpUpdatePlaySound
func webRdpUpdatePlaySound(context *C.rdpContext, playSound *C.cPlaySoundUpdate) C.BOOL {
	return C.TRUE
}

//export webRdpUpdateDesktopResize
func webRdpUpdateDesktopResize(context *C.rdpContext) C.BOOL {
	return C.TRUE
}

//export webRdpUpdateSetKeyboardIndicators
func webRdpUpdateSetKeyboardIndicators(context *C.rdpContext, ledFlags C.UINT16) C.BOOL {
	return C.TRUE
}

//export  webRdpUpdateSetKeyboardImeStatus
func webRdpUpdateSetKeyboardImeStatus(context *C.rdpContext, imeId C.UINT16, imeState C.UINT32, imeConvMode C.UINT32) C.BOOL {
	return C.TRUE
}

//export webRdpDstBlt
func webRdpDstBlt(context *C.rdpContext, dstblt *C.cDstbltOrder) C.BOOL {
	fmt.Println("webRdpDstBlt")
	return C.TRUE
}

//export webRdpOpaqueRect
func webRdpOpaqueRect(context *C.rdpContext, opaque_rect *C.cOpaqueRectOrder) C.BOOL {
	fmt.Println("webRdpOpaqueRect")
	return C.TRUE
}
