package freerdp

/*
#cgo CFLAGS: -I /opt/homebrew/Cellar/freerdp/2.9.0/include/freerdp2 -I /opt/homebrew/Cellar/freerdp/2.9.0/include/winpr2/
#cgo LDFLAGS: -L /opt/homebrew/Cellar/freerdp/2.9.0/lib  -lfreerdp2 -lfreerdp-client2 -lwinpr2

#if __APPLE__
#include <TargetConditionals.h>
#endif

#include "freerdp/freerdp.h"
#include "freerdp/client.h"

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
	ep->ContextSize = sizeof(rdpContext);
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
*/
import "C"
import (
	"fmt"
	"github.com/cihub/seelog"
)

type Client struct {
	ctx *C.rdpContext
}

func NewClient(addr, username, password string) *Client {
	seelog.Info("")

	var ep C.RDP_CLIENT_ENTRY_POINTS
	C.init_entry_points(&ep)

	ctx := C.freerdp_client_context_new(&ep)
	settings := ctx.instance.settings
	settings.ServerHostname = C.CString(addr)
	settings.Username = C.CString(username)
	settings.Password = C.CString(password)
	settings.IgnoreCertificate = C.TRUE

	instance := ctx.instance
	C.init_instance(instance)

	return &Client{ctx}
}

func (c *Client) Connect() error {
	context := c.ctx
	if C.freerdp_connect(context.instance) != C.TRUE {
		return fmt.Errorf("freerdp_connect return false")
	}
	return nil
}

func (c *Client) Start() error {
	context := c.ctx

	handlers := make([]C.HANDLE, 64)

	for C.freerdp_shall_disconnect(context.instance) != C.TRUE {
		r := C.freerdp_get_event_handles(context, &handlers[0], 64)
		if r == 0 {
			seelog.Info("freerdp_get_event_handles failed")
			break
		}
		if w := C.WaitForMultipleObjects(r, &handlers[0], C.FALSE, 250); w == C.WAIT_FAILED {
			seelog.Info("WaitForMultipleObjects failed: %v", w)
			break
		} else if w == C.WAIT_TIMEOUT {
			continue
		}
		if C.freerdp_check_event_handles(context) == C.FALSE {
			seelog.Info("freerdp_check_event_handles failed")
			break
		}
	}
	return nil
}

func (c *Client) DisConnect() {
	context := c.ctx
	C.freerdp_disconnect(context.instance)
}
