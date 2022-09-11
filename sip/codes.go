package sip

type ResponseCode int

const (
	Ok ResponseCode = 200

	Trying               ResponseCode = 100
	Ringing              ResponseCode = 180
	CallIsBeingForwarded ResponseCode = 181 //Call Is Being Forwarded
	Queued               ResponseCode = 182
	SessionProgress      ResponseCode = 183 //Session Progress

	MultipleChoices    ResponseCode = 300 //Multiple Choices
	MovedPermanently   ResponseCode = 301 //Moved Permanently
	MovedTemporarily   ResponseCode = 302 //Moved Temporarily
	UseProxy           ResponseCode = 305 //Use Proxy
	AlternativeService ResponseCode = 380 //Alternative Service

	BadRequest                  ResponseCode = 400 //Bad Request
	Unauthorized                ResponseCode = 401 //Unauthorized
	PaymentRequired             ResponseCode = 402 //Payment Required
	Forbidden                   ResponseCode = 403 //Forbidden
	NotFound                    ResponseCode = 404 //Not Found
	MethodNotAllowed            ResponseCode = 405 //Method Not Allowed
	NotAcceptable               ResponseCode = 406 //Not Acceptable
	ProxyAuthenticationRequired ResponseCode = 407 //Proxy Authentication Required
	RequestTimeout              ResponseCode = 408 //Request Timeout
	Gone                        ResponseCode = 410 //Gone
	RequestEntityTooLarge       ResponseCode = 413 //Request Entity Too Large
	RequestURITooLarge          ResponseCode = 414 //Request-URI Too Large
	UnsupportedMediaType        ResponseCode = 415 //Unsupported Media Type
	UnsupportedURIScheme        ResponseCode = 416 //Unsupported URI Scheme
	BadExtension                ResponseCode = 420 //Bad Extension
	ExtensionRequired           ResponseCode = 421 //Extension Required
	IntervalTooBrief            ResponseCode = 423 //Interval Too Brief
	TemporarilyNotAvailable     ResponseCode = 480 //Temporarily not available
	CallLegDoesNotExist         ResponseCode = 481 //Call Leg/Transaction Does Not Exist
	LoopDetected                ResponseCode = 482 //Loop Detected
	TooManyHops                 ResponseCode = 483 //Too Many Hops
	AddressIncomplete           ResponseCode = 484 //Address Incomplete
	Ambiguous                   ResponseCode = 485 //Ambiguous
	BusyHere                    ResponseCode = 486 //Busy Here
	RequestTerminated           ResponseCode = 487 //Request Terminated
	NotAcceptableHere           ResponseCode = 488 //Not Acceptable Here
	RequestPending              ResponseCode = 491 //Request Pending
	Undecipherable              ResponseCode = 493 //Undecipherable

	InternalServerError    ResponseCode = 500 //Internal Server Error
	NotImplemented         ResponseCode = 501 //Not Implemented
	BadGateway             ResponseCode = 502 //Bad Gateway
	ServiceUnavailable     ResponseCode = 503 //Service Unavailable
	ServerTimeout          ResponseCode = 504 //Server Time-out
	SIPVersionNotSupported ResponseCode = 505 //SIP Version not supported
	MessageTooLarge        ResponseCode = 513 //Message Too Large

	BusyEverywhere       ResponseCode = 600 //Busy Everywhere
	Decline              ResponseCode = 603 //Decline
	DoesNotExistAnywhere ResponseCode = 604 //Does not exist anywhere
	NotAcceptable606     ResponseCode = 606 //Not Acceptable
)

var ResponseCodes map[int]string = map[int]string{
	200: "OK",
	100: "Trying",
	180: "Ringing",
	181: "Call Is Being Forwarded",
	182: "Queued",
	183: "Session Progress",

	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Moved Temporarily",
	305: "Use Proxy",
	380: "Alternative Service",

	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	410: "Gone",
	413: "Request Entity Too Large",
	414: "Request-URI Too Large",
	415: "Unsupported Media Type",
	416: "Unsupported URI Scheme",
	420: "Bad Extension",
	421: "Extension Required",
	423: "Interval Too Brief",
	480: "Temporarily not available",
	481: "Call Leg/Transaction Does Not Exist",
	482: "Loop Detected",
	483: "Too Many Hops",
	484: "Address Incomplete",
	485: "Ambiguous",
	486: "Busy Here",
	487: "Request Terminated",
	488: "Not Acceptable Here",
	491: "Request Pending",
	493: "Undecipherable",

	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Server Time-out",
	505: "SIP Version not supported",
	513: "Message Too Large",

	600: "Busy Everywhere",
	603: "Decline",
	604: "Does not exist anywhere",
	606: "Not Acceptable",
}
