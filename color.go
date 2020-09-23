package regip

import (
	"github.com/gookit/color"
)

var (
	// Areas - Server
	CLR_server    = color.Cyan
	CLR_session   = color.Blue
	CLR_readloop  = color.Yellow
	CLR_writeloop = color.Red
	CLR_chainloop = color.White

	CLR_api = color.Green

	// Areas - Client
	CLR_cli = color.White

	// Types
	CLR_mt   = color.Magenta
	CLR_time = color.Yellow

	// Areas - Database
	CLR_db           = color.White
	CLR_fullindex    = color.Yellow
	CLR_indexrecords = color.Blue
)
