package freerdp

/*
#cgo darwin CFLAGS: -I /opt/homebrew/Cellar/freerdp/2.9.0/include/freerdp2 -I /opt/homebrew/Cellar/freerdp/2.9.0/include/winpr2/
#cgo darwin LDFLAGS: -L /opt/homebrew/Cellar/freerdp/2.9.0/lib -lfreerdp2 -lfreerdp-client2 -lwinpr2
#cgo linux CFLAGS: -I /usr/include/winpr2/ -I /usr/include/freerdp2/
#cgo linux LDFLAGS: -lfreerdp2 -lfreerdp-client2 -lwinpr2

#if __APPLE__
#include <TargetConditionals.h>
#endif

#include "freerdp/freerdp.h"
#include "freerdp/client.h"

typedef struct {
	rdpContext context;
	void *priv;
} webRdpContent;

static uintptr_t webrdp_content_get(rdpContext *content)
{
	webRdpContent *wctx = (webRdpContent*)content;
	return (uintptr_t)wctx->priv;
}

static void webrdp_content_set(rdpContext *content, uintptr_t priv)
{
	webRdpContent *wctx = (webRdpContent*)content;
	wctx->priv = (void*)priv;
}

extern BOOL webRdpClientGlobalInit();
extern void webRdpClientGlobalUnInit();
extern BOOL webRdpClientNew(freerdp *instance, rdpContext* context);
extern void webRdpClientFree(freerdp *instance, rdpContext* context);
extern int webRdpClientStart(rdpContext *context);
extern int webRdpClientStop(rdpContext *context);

static void init_entry_points(RDP_CLIENT_ENTRY_POINTS *ep)
{
	ep->Size = sizeof(RDP_CLIENT_ENTRY_POINTS);
	ep->Version = RDP_CLIENT_INTERFACE_VERSION;
	ep->ContextSize = sizeof(webRdpContent);
	ep->GlobalInit = webRdpClientGlobalInit;
	ep->GlobalUninit = webRdpClientGlobalUnInit;
	ep->ClientNew = webRdpClientNew;
	ep->ClientFree = webRdpClientFree;
	ep->ClientStart = webRdpClientStart;
	ep->ClientStop = webRdpClientStop;
}

extern BOOL webRdpClientPreConnect(freerdp* instance);
extern BOOL webRdpClientPostConnect(freerdp* instance);
extern BOOL webRdpClientAuthenticate(freerdp* instance, char** username, char** password, char** domain);

static void init_instance(freerdp *instance)
{
	instance->PreConnect = webRdpClientPreConnect;
	instance->PostConnect = webRdpClientPostConnect;
	instance->Authenticate = webRdpClientAuthenticate;
}

static void send_mouse_event(freerdp *instance, UINT16 flags, int x, int y)
{
	rdpInput *input = instance->input;
	input->MouseEvent(input, flags, x, y);
}

static void send_keyboard_event(freerdp *instance, UINT16 flags, int k)
{
	rdpInput *input = instance->input;
	input->KeyboardEvent(input, flags, k);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Client struct {
	context  *C.rdpContext
	dataChan chan *Message
	isClosed bool
}

func NewClient(addr, username, password string) *Client {
	var ep C.RDP_CLIENT_ENTRY_POINTS
	C.init_entry_points(&ep)

	context := C.freerdp_client_context_new(&ep)
	settings := context.instance.settings
	settings.ServerHostname = C.CString(addr)
	settings.Username = C.CString(username)
	settings.Password = C.CString(password)
	settings.IgnoreCertificate = C.TRUE

	instance := context.instance
	C.init_instance(instance)

	c := &Client{
		context:  context,
		dataChan: make(chan *Message, 8192),
	}

	return c.setClientToContent(context)
}

func (c *Client) Connect() error {
	context := c.context
	if C.freerdp_connect(context.instance) != C.TRUE {
		return fmt.Errorf("freerdp_connect return false")
	}
	return nil
}

func (c *Client) setClientToContent(context *C.rdpContext) *Client {
	C.webrdp_content_set(context, C.uintptr_t(uintptr(unsafe.Pointer(c))))
	return c
}

func getClientFromContent(context *C.rdpContext) *Client {
	ptr := C.webrdp_content_get(context)
	return (*Client)(unsafe.Pointer(uintptr(ptr)))
}

func (c *Client) Start() error {
	context := c.context
	handlers := make([]C.HANDLE, 64)

	for C.freerdp_shall_disconnect(context.instance) != C.TRUE {
		r := C.freerdp_get_event_handles(context, &handlers[0], 64)
		if r == 0 {
			break
		}
		if w := C.WaitForMultipleObjects(r, &handlers[0], C.FALSE, 250); w == C.WAIT_FAILED {
			break
		} else if w == C.WAIT_TIMEOUT {
			continue
		}
		if C.freerdp_check_event_handles(context) == C.FALSE {
			break
		}
	}
	return nil
}

func (c *Client) DisConnect() {
	fmt.Println("---- ")

	if c.isClosed {
		return
	}

	c.isClosed = true
	close(c.dataChan)

	context := c.context
	C.freerdp_disconnect(context.instance)

}

func (c *Client) Data() (*Message, bool) {
	data, ok := <-c.dataChan
	return data, ok
}

// 鼠标事件
func (c *Client) ProcessMouseEvent(mouse *Mouse) {
	switch mouse.Type {
	case "mousemove":
		C.send_mouse_event(c.context.instance, C.PTR_FLAGS_MOVE, C.int(mouse.X), C.int(mouse.Y))
	case "mousedown":
		switch mouse.Btn {
		case 0: // 左
			C.send_mouse_event(c.context.instance, C.PTR_FLAGS_DOWN|C.PTR_FLAGS_BUTTON1, C.int(mouse.X), C.int(mouse.Y))
		case 1: // 中
			C.send_mouse_event(c.context.instance, C.PTR_FLAGS_DOWN|C.PTR_FLAGS_BUTTON3, C.int(mouse.X), C.int(mouse.Y))
		case 2: // 右
			C.send_mouse_event(c.context.instance, C.PTR_FLAGS_DOWN|C.PTR_FLAGS_BUTTON2, C.int(mouse.X), C.int(mouse.Y))
		}
	case "mouseup":
		switch mouse.Btn {
		case 0: // 左
			C.send_mouse_event(c.context.instance, C.PTR_FLAGS_BUTTON1, C.int(mouse.X), C.int(mouse.Y))
		case 1: // 中
			C.send_mouse_event(c.context.instance, C.PTR_FLAGS_BUTTON3, C.int(mouse.X), C.int(mouse.Y))
		case 2: // 右
			C.send_mouse_event(c.context.instance, C.PTR_FLAGS_BUTTON2, C.int(mouse.X), C.int(mouse.Y))
		}
	}
}

var keyMap = map[int]int{
	8:   0xE,    // Backspace
	9:   0xF,    // Tab
	13:  0x1C,   // Enter
	16:  0x2A,   // ShiftLeft
	17:  0x1D,   // ControlLeft
	18:  0x38,   // AltLeft
	20:  0x3A,   // CapsLock
	27:  0x1,    // Escape
	32:  0x39,   // Space
	37:  0xE04B, // ArrowLeft
	38:  0xE048, // ArrowUp
	39:  0xE04D, // ArrowRight
	40:  0xE050, // ArrowDown
	48:  0xB,    // Digit0
	49:  0x2,    // Digit1
	50:  0x3,    // Digit2
	51:  0x4,    // Digit3
	52:  0x5,    // Digit4
	53:  0x6,    // Digit5
	54:  0x7,    // Digit6
	55:  0x8,    // Digit7
	56:  0x9,    // Digit8
	57:  0xA,    // Digit9
	59:  0x27,   // Semicolon
	61:  0xD,    // Equal
	65:  0x1E,   // KeyA
	66:  0x30,   // KeyB
	67:  0x2E,   // KeyC
	68:  0x20,   // KeyD
	69:  0x12,   // KeyE
	70:  0x21,   // KeyF
	71:  0x22,   // KeyG
	72:  0x23,   // KeyH
	73:  0x17,   // KeyI
	74:  0x24,   // KeyJ
	75:  0x25,   // KeyK
	76:  0x26,   // KeyL
	77:  0x32,   // KeyM
	78:  0x31,   // KeyN
	79:  0x18,   // KeyO
	80:  0x19,   // KeyP
	81:  0x10,   // KeyQ
	82:  0x13,   // KeyR
	83:  0x1F,   // KeyS
	84:  0x14,   // KeyT
	85:  0x16,   // KeyU
	86:  0x2F,   // KeyV
	87:  0x11,   // KeyW
	88:  0x2D,   // KeyX
	89:  0x15,   // KeyY
	90:  0x2C,   // KeyZ
	91:  0xE05B, // OSLeft
	93:  0xE05D, // ContextMenu
	96:  0x52,   // Numpad0
	97:  0x4F,   // Numpad1
	98:  0x50,   // Numpad2
	99:  0x51,   // Numpad3
	100: 0x4B,   // Numpad4
	101: 0x4C,   // Numpad5
	102: 0x4D,   // Numpad6
	103: 0x47,   // Numpad7
	104: 0x48,   // Numpad8
	105: 0x49,   // Numpad9
	106: 0x37,   // NumpadMultiply
	107: 0x4E,   // NumpadAdd
	109: 0x4A,   // NumpadSubtract
	110: 0x53,   // NumpadDecimal
	111: 0xE035, // NumpadDivide
	112: 0x3B,   // F1
	113: 0x3C,   // F2
	114: 0x3D,   // F3
	115: 0x3E,   // F4
	116: 0x3F,   // F5
	117: 0x40,   // F6
	118: 0x41,   // F7
	119: 0x42,   // F8
	120: 0x43,   // F9
	121: 0x44,   // F10
	122: 0x57,   // F11
	123: 0x58,   // F12
	144: 0xE045, // NumLock
	173: 0xC,    // Minus
	188: 0x33,   // Comma
	190: 0x34,   // Period
	191: 0x35,   // Slash
	192: 0x29,   // Backquote
	219: 0x1A,   // BracketLeft
	220: 0x2B,   // Backslash
	221: 0x1B,   // BracketRight
	222: 0x28,   // Quote
}

func (c *Client) ProcessKeyboardEvent(keyboard *Keyboard) {
	switch keyboard.Type {
	case "keydown":
		C.send_keyboard_event(c.context.instance, C.KBD_FLAGS_DOWN, C.int(keyMap[keyboard.K]))
	case "keyup":
		C.send_keyboard_event(c.context.instance, C.KBD_FLAGS_RELEASE, C.int(keyMap[keyboard.K]))
	}
}
