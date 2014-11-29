package keyring

/*
#include <CoreFoundation/CoreFoundation.h>
#include <Security/Security.h>
#include <CoreServices/CoreServices.h>

#cgo LDFLAGS: -framework CoreFoundation -framework Security

static char*
cstring(CFStringRef s)
{
	char *p;
	int n;
	n = CFStringGetLength(s)*8;
	p = malloc(n);
	CFStringGetCString(s, p, n, kCFStringEncodingUTF8);
	return p;
}

void
getpasswd(char *service, char *user, char **passwd, char **error)
{
	OSStatus status;
	CFStringRef str;
	UInt32 len;
	void *data;

	*passwd = NULL;
	*error = NULL;
	status = SecKeychainFindInternetPassword(
		NULL,                     // default keychain
		strlen(service), service, // service name
		0, NULL,                  // security domain
		strlen(user), user,       // account name
		0, NULL,                  // path
		0,                        // port
		0,                        // protocol type
		kSecAuthenticationTypeDefault,
		&len,
		&data,
		NULL
	);
	if(status != 0){
		str = SecCopyErrorMessageString(status, NULL);
		*error = cstring(str);
		CFRelease(str);
		return;
	}

	*passwd = malloc(len+1);
	memmove(*passwd, data, len);
	(*passwd)[len] = '\0';
}

void
putpasswd(char *service, char *user, char *passwd, char **error)
{
	OSStatus status;
	CFStringRef str;

	status = SecKeychainAddInternetPassword(
		NULL,                     // default keychain
		strlen(service), service, // service name
		0, NULL,                  // security domain
		strlen(user), user,       // account name
		0, NULL,                  // path
		0,                        // port
		0,                        // protocol type
		kSecAuthenticationTypeDefault,
		strlen(passwd), passwd,   // password
		NULL
	);
	if(status != 0){
		str = SecCopyErrorMessageString(status, NULL);
		*error = cstring(str);
		CFRelease(str);
	}
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

func init() { keyring["darwin"] = &keychain{} }

type keychain struct{}

func (k *keychain) Get(service, username string) (string, error) {
	cservice := C.CString(service)
	cuser := C.CString(username)
	defer C.free(unsafe.Pointer(cservice))
	defer C.free(unsafe.Pointer(cuser))

	var cpasswd, cerror *C.char
	C.getpasswd(cservice, cuser, &cpasswd, &cerror)
	defer C.free(unsafe.Pointer(cpasswd))
	defer C.free(unsafe.Pointer(cerror))
	if cerror != nil {
		return "", errors.New(C.GoString(cerror))
	}
	return C.GoString(cpasswd), nil
}

func (k *keychain) Set(service, username, password string) error {
	cservice := C.CString(service)
	cuser := C.CString(username)
	cpasswd := C.CString(password)
	defer C.free(unsafe.Pointer(cservice))
	defer C.free(unsafe.Pointer(cuser))
	defer C.free(unsafe.Pointer(cpasswd))

	var cerror *C.char
	defer C.free(unsafe.Pointer(cerror))
	C.putpasswd(cservice, cuser, cpasswd, &cerror)
	if cerror != nil {
		return errors.New(C.GoString(cerror))
	}
	return nil
}

func (k *keychain) Delete(service, username string) error {
	panic("not implemented")
}
