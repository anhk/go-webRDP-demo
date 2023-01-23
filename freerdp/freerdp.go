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
	case "mosedown":
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
