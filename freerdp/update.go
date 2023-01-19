package freerdp

/*
#if __APPLE__
#include <TargetConditionals.h>
#endif

#include "freerdp/freerdp.h"

typedef const PLAY_SOUND_UPDATE cPsu;

extern BOOL webRdpUpdateBeginPaint(rdpContext* context);
extern BOOL webRdpUpdateEndPaint(rdpContext* context);
extern BOOL webRdpUpdatePlaySound(rdpContext* context, cPsu* play_sound);
extern BOOL webRdpUpdateDesktopResize(rdpContext* context);
extern BOOL webRdpUpdateSetKeyboardIndicators(rdpContext* context, UINT16 led_flags);
extern BOOL webRdpUpdateSetKeyboardImeStatus(rdpContext* context, UINT16 imeId, UINT32 imeState, UINT32 imeConvMode);

static void init_update_callbacks(rdpUpdate *update)
{
    update->BeginPaint = webRdpUpdateBeginPaint;
    update->EndPaint = webRdpUpdateEndPaint;
    update->PlaySound = webRdpUpdatePlaySound;
    update->DesktopResize = webRdpUpdateDesktopResize;
    update->SetKeyboardIndicators = webRdpUpdateSetKeyboardIndicators;
    update->SetKeyboardImeStatus = webRdpUpdateSetKeyboardImeStatus;
}
*/
import "C"

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
func webRdpUpdatePlaySound(context *C.rdpContext, playSound *C.cPsu) C.BOOL {
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
