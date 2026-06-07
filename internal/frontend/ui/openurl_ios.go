//go:build ios

package ui

/*
#cgo CFLAGS: -x objective-c -fblocks
#cgo LDFLAGS: -framework Foundation -framework UIKit

#import <dispatch/dispatch.h>
#import <Foundation/Foundation.h>
#import <stdlib.h>
#import <UIKit/UIKit.h>

static void passgoOpenURL(const char *rawURL) {
	@autoreleasepool {
		NSString *urlString = [NSString stringWithUTF8String:rawURL];
		NSURL *url = [NSURL URLWithString:urlString];
		if (url == nil) {
			return;
		}
		dispatch_async(dispatch_get_main_queue(), ^{
			[[UIApplication sharedApplication] openURL:url options:@{} completionHandler:nil];
		});
	}
}
*/
import "C"
import "unsafe"

func OpenURL(url string) error {
	url = NormalizeURL(url)
	if url == "" {
		return nil
	}

	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	C.passgoOpenURL(cURL)
	return nil
}
