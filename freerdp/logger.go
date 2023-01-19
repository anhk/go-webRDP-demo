package freerdp

import (
	"github.com/cihub/seelog"
)

var SeeLogConfig = `
<seelog type="asynctimer" asyncinterval="1000000" minlevel="debug" maxlevel="error">
    <outputs formatid="main">
        <console/>
    </outputs>
    <formats>
        <format id="main" format="[%UTCDate %UTCTime][%LEVEL][%RelFile:%Line][%FuncShort] %Msg%n"/>
    </formats>
</seelog>
`

func init() {
	if logger, err := seelog.LoggerFromConfigAsString(SeeLogConfig); err == nil {
		_ = seelog.ReplaceLogger(logger)
	}
}
