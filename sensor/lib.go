// Package sensor provides a minimal, cgo-free wrapper around the Android
// NDK sensor API. It dynamically links against libandroid.so at runtime via
// purego, exposing just the subset of the API that this program needs.
package sensor

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// C function signatures (matching the arm64 Android C ABI). Pointers are
// represented as uintptr; const char* results are converted to Go strings
// by purego.
var (
	fnGetInstance    func() uintptr                                    // ASensorManager* ASensorManager_getInstance(void)
	fnGetSensorList  func(manager uintptr, list *unsafe.Pointer) int32 // int ASensorManager_getSensorList(ASensorManager*, ASensorList*)
	fnGetName        func(sensor uintptr) string                       // const char* ASensor_getName(const ASensor*)
	fnGetVendor      func(sensor uintptr) string                       // const char* ASensor_getVendor(const ASensor*)
	fnGetType        func(sensor uintptr) int32                        // int ASensor_getType(const ASensor*)
	fnGetStringType  func(sensor uintptr) string                       // const char* ASensor_getStringType(const ASensor*)
	fnGetMinDelay    func(sensor uintptr) int32                        // int ASensor_getMinDelay(const ASensor*)
	fnGetResolution  func(sensor uintptr) float32                      // float ASensor_getResolution(const ASensor*)
	fnReportingMode  func(sensor uintptr) int32                        // int ASensor_getReportingMode(const ASensor*)
	fnFifoMax        func(sensor uintptr) int32                        // int ASensor_getFifoMaxEventCount(const ASensor*)
	fnFifoReserved   func(sensor uintptr) int32                        // int ASensor_getFifoReservedEventCount(const ASensor*)
	fnGetHandle      func(sensor uintptr) int32                        // int ASensor_getHandle(const ASensor*)
	fnDirectRate     func(sensor uintptr) int32                        // int ASensor_getHighestDirectReportRateLevel(const ASensor*)
	fnIsWakeUpSensor func(sensor uintptr) bool                         // bool ASensor_isWakeUpSensor(const ASensor*)

	fnCreateEventQueue  func(manager, looper uintptr, ident int32, callback, data uintptr) uintptr // ASensorEventQueue* ASensorManager_createEventQueue(...)
	fnDestroyEventQueue func(manager, queue uintptr) int32                                         // int ASensorManager_destroyEventQueue(...)
	fnEnableSensor      func(queue, sensor uintptr) int32                                          // int ASensorEventQueue_enableSensor(...)
	fnDisableSensor     func(queue, sensor uintptr) int32                                          // int ASensorEventQueue_disableSensor(...)
	fnHasEvents         func(queue uintptr) int32                                                  // int ASensorEventQueue_hasEvents(...)
	fnGetEvents         func(queue uintptr, events *cEvent, count uintptr) int64                   // ssize_t ASensorEventQueue_getEvents(...)

	fnALooperPrepare func(opts int32) uintptr // ALooper* ALooper_prepare(int opts)
	fnALooperRelease func(looper uintptr)     // void ALooper_release(ALooper* looper)
)

// cEvent mirrors the C ASensorEvent struct (arm64 layout, 104 bytes).
type cEvent struct {
	Version   int32
	Sensor    int32
	Type      int32
	Reserved0 int32
	Timestamp int64
	Data      [16]float32
	Flags     uint32
	Reserved1 [3]int32
}

var (
	libHandle uintptr
	loadOnce  sync.Once
	loadErr   error
)

// openLibandroid loads the Android system library. The system path is tried
// first because some environments (e.g. Termux) ship a reduced libandroid.so
// shim that shadows the real one when resolving by bare name.
func openLibandroid() (uintptr, error) {
	systemPath := "/system/lib/libandroid.so"
	if unsafe.Sizeof(uintptr(0)) == 8 {
		systemPath = "/system/lib64/libandroid.so"
	}

	var lastErr error
	for _, path := range []string{systemPath, "libandroid.so"} {
		handle, err := purego.Dlopen(path, purego.RTLD_NOW)
		if err == nil {
			return handle, nil
		}
		lastErr = err
	}
	return 0, lastErr
}

// load resolves every native symbol once and caches the result.
func load() error {
	loadOnce.Do(func() {
		libHandle, loadErr = openLibandroid()
		if loadErr != nil {
			return
		}

		bindings := []struct {
			fptr any
			name string
		}{
			{&fnGetInstance, "ASensorManager_getInstance"},
			{&fnGetSensorList, "ASensorManager_getSensorList"},
			{&fnGetName, "ASensor_getName"},
			{&fnGetVendor, "ASensor_getVendor"},
			{&fnGetType, "ASensor_getType"},
			{&fnGetStringType, "ASensor_getStringType"},
			{&fnGetMinDelay, "ASensor_getMinDelay"},
			{&fnGetResolution, "ASensor_getResolution"},
			{&fnReportingMode, "ASensor_getReportingMode"},
			{&fnFifoMax, "ASensor_getFifoMaxEventCount"},
			{&fnFifoReserved, "ASensor_getFifoReservedEventCount"},
			{&fnGetHandle, "ASensor_getHandle"},
			{&fnDirectRate, "ASensor_getHighestDirectReportRateLevel"},
			{&fnIsWakeUpSensor, "ASensor_isWakeUpSensor"},
			{&fnCreateEventQueue, "ASensorManager_createEventQueue"},
			{&fnDestroyEventQueue, "ASensorManager_destroyEventQueue"},
			{&fnEnableSensor, "ASensorEventQueue_enableSensor"},
			{&fnDisableSensor, "ASensorEventQueue_disableSensor"},
			{&fnHasEvents, "ASensorEventQueue_hasEvents"},
			{&fnGetEvents, "ASensorEventQueue_getEvents"},
			{&fnALooperPrepare, "ALooper_prepare"},
			{&fnALooperRelease, "ALooper_release"},
		}

		for _, b := range bindings {
			addr, err := purego.Dlsym(libHandle, b.name)
			if err != nil {
				loadErr = fmt.Errorf("resolve %s: %w", b.name, err)
				return
			}
			purego.RegisterFunc(b.fptr, addr)
		}
	})
	return loadErr
}
