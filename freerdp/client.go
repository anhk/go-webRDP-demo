package freerdp

/*
#if __APPLE__
#include <TargetConditionals.h>
#endif

#include "freerdp/freerdp.h"
#include "freerdp/cache/cache.h"
#include "freerdp/gdi/gdi.h"
*/
import "C"

//export webRdpClientGlobalInit
func webRdpClientGlobalInit() C.BOOL {
	return C.TRUE
}

//export webRdpClientGlobalUnInit
func webRdpClientGlobalUnInit() {
}

//export webRdpClientNew
func webRdpClientNew(instance *C.freerdp, context *C.rdpContext) C.BOOL {
	return C.TRUE
}

//export webRdpClientFree
func webRdpClientFree(instance *C.freerdp, context *C.rdpContext) {
}

//export webRdpClientStart
func webRdpClientStart(context *C.rdpContext) C.int {
	return 0
}

//export webRdpClientStop
func webRdpClientStop(context *C.rdpContext) C.int {
	return 0
}

//export webRdpClientPreConnect
func webRdpClientPreConnect(instance *C.freerdp) C.BOOL {
	settings := instance.settings
	settings.OsMajorType = C.OSMAJORTYPE_UNIX
	settings.OsMinorType = C.OSMINORTYPE_NATIVE_XSERVER
	settings.FastPathOutput = 1
	settings.ColorDepth = 32
	settings.BitmapCacheV3Enabled = 0
	settings.BitmapCachePersistEnabled = 0
	settings.GlyphSupportLevel = C.GLYPH_SUPPORT_NONE

	if instance.context.cache == nil {
		instance.context.cache = C.cache_new(instance.settings)
	}
	return C.TRUE
}

//export webRdpClientPostConnect
func webRdpClientPostConnect(instance *C.freerdp) C.BOOL {
	if C.gdi_init(instance, C.PIXEL_FORMAT_XRGB32) != C.TRUE {
		return C.FALSE
	}

	registerGraphics(instance.context.graphics)
	registerUpdateCallbacks(instance.update)

	return C.TRUE
}

//export webRdpClientAuthenticate
func webRdpClientAuthenticate(instance *C.freerdp, username **C.char,
	password **C.char, domain **C.char) C.BOOL {

	return C.TRUE
}
