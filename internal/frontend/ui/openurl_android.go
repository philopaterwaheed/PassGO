//go:build android && cgo

package ui

/*
#include <jni.h>
#include <stdlib.h>

static JavaVM *passgoVM = NULL;

JNIEXPORT jint JNICALL JNI_OnLoad(JavaVM *vm, void *reserved) {
	passgoVM = vm;
	return JNI_VERSION_1_6;
}

static int passgoException(JNIEnv *env, int code) {
	if ((*env)->ExceptionCheck(env)) {
		(*env)->ExceptionClear(env);
		return code;
	}
	return 0;
}

static int passgoOpenURLAndroid(const char *rawURL) {
	JavaVM *vm = passgoVM;
	if (vm == NULL) {
		return -1;
	}

	JNIEnv *env = NULL;
	int detach = 0;
	jint envStatus = (*vm)->GetEnv(vm, (void **)&env, JNI_VERSION_1_6);
	if (envStatus == JNI_EDETACHED) {
		if ((*vm)->AttachCurrentThread(vm, (void **)&env, NULL) != JNI_OK) {
			return -2;
		}
		detach = 1;
	} else if (envStatus != JNI_OK || env == NULL) {
		return -3;
	}

	jclass activityThreadCls = (*env)->FindClass(env, "android/app/ActivityThread");
	if (passgoException(env, -4) != 0 || activityThreadCls == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -4;
	}

	jmethodID currentApplication = (*env)->GetStaticMethodID(env, activityThreadCls, "currentApplication", "()Landroid/app/Application;");
	if (passgoException(env, -5) != 0 || currentApplication == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -5;
	}

	jobject context = (*env)->CallStaticObjectMethod(env, activityThreadCls, currentApplication);
	if (passgoException(env, -6) != 0 || context == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -6;
	}

	jclass uriCls = (*env)->FindClass(env, "android/net/Uri");
	if (passgoException(env, -7) != 0 || uriCls == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -7;
	}

	jmethodID parseURI = (*env)->GetStaticMethodID(env, uriCls, "parse", "(Ljava/lang/String;)Landroid/net/Uri;");
	if (passgoException(env, -8) != 0 || parseURI == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -8;
	}

	jstring urlString = (*env)->NewStringUTF(env, rawURL);
	if (passgoException(env, -9) != 0 || urlString == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -9;
	}

	jobject uri = (*env)->CallStaticObjectMethod(env, uriCls, parseURI, urlString);
	if (passgoException(env, -10) != 0 || uri == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -10;
	}

	jclass intentCls = (*env)->FindClass(env, "android/content/Intent");
	if (passgoException(env, -11) != 0 || intentCls == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -11;
	}

	jmethodID newIntent = (*env)->GetMethodID(env, intentCls, "<init>", "(Ljava/lang/String;Landroid/net/Uri;)V");
	if (passgoException(env, -12) != 0 || newIntent == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -12;
	}

	jstring actionView = (*env)->NewStringUTF(env, "android.intent.action.VIEW");
	if (passgoException(env, -13) != 0 || actionView == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -13;
	}

	jobject intent = (*env)->NewObject(env, intentCls, newIntent, actionView, uri);
	if (passgoException(env, -14) != 0 || intent == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -14;
	}

	jmethodID addFlags = (*env)->GetMethodID(env, intentCls, "addFlags", "(I)Landroid/content/Intent;");
	if (passgoException(env, -15) != 0 || addFlags == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -15;
	}

	(*env)->CallObjectMethod(env, intent, addFlags, 0x10000000); // Intent.FLAG_ACTIVITY_NEW_TASK.
	if (passgoException(env, -16) != 0) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -16;
	}

	jclass contextCls = (*env)->GetObjectClass(env, context);
	if (passgoException(env, -17) != 0 || contextCls == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -17;
	}

	jmethodID startActivity = (*env)->GetMethodID(env, contextCls, "startActivity", "(Landroid/content/Intent;)V");
	if (passgoException(env, -18) != 0 || startActivity == NULL) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -18;
	}

	(*env)->CallVoidMethod(env, context, startActivity, intent);
	if (passgoException(env, -19) != 0) {
		if (detach) (*vm)->DetachCurrentThread(vm);
		return -19;
	}

	if (detach) {
		(*vm)->DetachCurrentThread(vm);
	}
	return 0;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func OpenURL(url string) error {
	url = NormalizeURL(url)
	if url == "" {
		return nil
	}

	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	if code := C.passgoOpenURLAndroid(cURL); code != 0 {
		return fmt.Errorf("open android url: jni error %d", int(code))
	}
	return nil
}
