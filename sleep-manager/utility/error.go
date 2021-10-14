package utility

import "errors"

var (
	ErrorBedinTime error = errors.New("TimeError: you don't sleep in p.m 0 ~ p.m 9")
	ErrorVerify          = errors.New("VetifyRrror: failed to verify signature")
	ErrorInit            = errors.New("InitError: failed to create bot")
	ErrorReply           = errors.New("ReplyError: failed: can't reply message")
)
