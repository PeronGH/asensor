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
	fnGetMinDelay    func(sensor uintptr) int32                        // int ASensor_getMinDelay(const ASensor*)
	fnGetResolution  func(sensor uintptr) float32                      // float ASensor_getResolution(const ASensor*)
	fnIsWakeUpSensor func(sensor uintptr) bool                         // bool ASensor_isWakeUpSensor(const ASensor*)
)

var (
	libHandle uintptr
	loadOnce  sync.Once
	loadErr   error
)

// load resolves every native symbol once and caches the result.
func load() error {
	loadOnce.Do(func() {
		libHandle, loadErr = purego.Dlopen("libandroid.so", purego.RTLD_NOW)
		if loadErr != nil {
			loadErr = fmt.Errorf("dlopen libandroid.so: %w", loadErr)
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
			{&fnGetMinDelay, "ASensor_getMinDelay"},
			{&fnGetResolution, "ASensor_getResolution"},
			{&fnIsWakeUpSensor, "ASensor_isWakeUpSensor"},
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
