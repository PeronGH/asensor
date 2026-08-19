package sensor

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

// Type is an Android sensor type as returned by ASensor_getType.
type Type int32

const (
	Invalid                              Type = -1
	Accelerometer                        Type = 1
	MagneticField                        Type = 2
	Orientation                          Type = 3
	Gyroscope                            Type = 4
	Light                                Type = 5
	Pressure                             Type = 6
	Temperature                          Type = 7
	Proximity                            Type = 8
	Gravity                              Type = 9
	LinearAcceleration                   Type = 10
	RotationVector                       Type = 11
	RelativeHumidity                     Type = 12
	AmbientTemperature                   Type = 13
	MagneticFieldUncalibrated            Type = 14
	GameRotationVector                   Type = 15
	GyroscopeUncalibrated                Type = 16
	SignificantMotion                    Type = 17
	StepDetector                         Type = 18
	StepCounter                          Type = 19
	GeomagneticRotationVector            Type = 20
	HeartRate                            Type = 21
	TiltDetector                         Type = 22
	WakeGesture                          Type = 23
	GlanceGesture                        Type = 24
	PickUpGesture                        Type = 25
	WristTiltGesture                     Type = 26
	DeviceOrientation                    Type = 27
	Pose6dof                             Type = 28
	StationaryDetect                     Type = 29
	MotionDetect                         Type = 30
	HeartBeat                            Type = 31
	DynamicSensorMeta                    Type = 32
	AdditionalInfo                       Type = 33
	LowLatencyOffbodyDetect              Type = 34
	AccelerometerUncalibrated            Type = 35
	HingeAngle                           Type = 36
	HeadTracker                          Type = 37
	AccelerometerLimitedAxes             Type = 38
	GyroscopeLimitedAxes                 Type = 39
	AccelerometerLimitedAxesUncalibrated Type = 40
	GyroscopeLimitedAxesUncalibrated     Type = 41
	Heading                              Type = 42
)

// String returns the Android name of the sensor type.
func (t Type) String() string {
	switch t {
	case Invalid:
		return "Invalid"
	case Accelerometer:
		return "Accelerometer"
	case MagneticField:
		return "MagneticField"
	case Orientation:
		return "Orientation"
	case Gyroscope:
		return "Gyroscope"
	case Light:
		return "Light"
	case Pressure:
		return "Pressure"
	case Temperature:
		return "Temperature"
	case Proximity:
		return "Proximity"
	case Gravity:
		return "Gravity"
	case LinearAcceleration:
		return "LinearAcceleration"
	case RotationVector:
		return "RotationVector"
	case RelativeHumidity:
		return "RelativeHumidity"
	case AmbientTemperature:
		return "AmbientTemperature"
	case MagneticFieldUncalibrated:
		return "MagneticFieldUncalibrated"
	case GameRotationVector:
		return "GameRotationVector"
	case GyroscopeUncalibrated:
		return "GyroscopeUncalibrated"
	case SignificantMotion:
		return "SignificantMotion"
	case StepDetector:
		return "StepDetector"
	case StepCounter:
		return "StepCounter"
	case GeomagneticRotationVector:
		return "GeomagneticRotationVector"
	case HeartRate:
		return "HeartRate"
	case TiltDetector:
		return "TiltDetector"
	case WakeGesture:
		return "WakeGesture"
	case GlanceGesture:
		return "GlanceGesture"
	case PickUpGesture:
		return "PickUpGesture"
	case WristTiltGesture:
		return "WristTiltGesture"
	case DeviceOrientation:
		return "DeviceOrientation"
	case Pose6dof:
		return "Pose6dof"
	case StationaryDetect:
		return "StationaryDetect"
	case MotionDetect:
		return "MotionDetect"
	case HeartBeat:
		return "HeartBeat"
	case DynamicSensorMeta:
		return "DynamicSensorMeta"
	case AdditionalInfo:
		return "AdditionalInfo"
	case LowLatencyOffbodyDetect:
		return "LowLatencyOffbodyDetect"
	case AccelerometerUncalibrated:
		return "AccelerometerUncalibrated"
	case HingeAngle:
		return "HingeAngle"
	case HeadTracker:
		return "HeadTracker"
	case AccelerometerLimitedAxes:
		return "AccelerometerLimitedAxes"
	case GyroscopeLimitedAxes:
		return "GyroscopeLimitedAxes"
	case AccelerometerLimitedAxesUncalibrated:
		return "AccelerometerLimitedAxesUncalibrated"
	case GyroscopeLimitedAxesUncalibrated:
		return "GyroscopeLimitedAxesUncalibrated"
	case Heading:
		return "Heading"
	default:
		return fmt.Sprintf("Type(%d)", int(t))
	}
}

// ValueCount returns how many leading floats of an event's data are
// meaningful for this sensor type.
func (t Type) ValueCount() int {
	switch t {
	case Light, Pressure, Proximity, RelativeHumidity, AmbientTemperature, Temperature,
		HeartRate, HeartBeat, SignificantMotion, StepDetector, StepCounter,
		StationaryDetect, MotionDetect, TiltDetector, WakeGesture, GlanceGesture,
		PickUpGesture, WristTiltGesture, DeviceOrientation, HingeAngle,
		LowLatencyOffbodyDetect:
		return 1
	case MagneticFieldUncalibrated, GyroscopeUncalibrated, AccelerometerUncalibrated,
		AccelerometerLimitedAxes, GyroscopeLimitedAxes:
		return 6
	case AccelerometerLimitedAxesUncalibrated, GyroscopeLimitedAxesUncalibrated:
		return 9
	case RotationVector, GameRotationVector, GeomagneticRotationVector:
		return 4
	case Heading:
		return 2
	default:
		return 3
	}
}

// ReportingMode is a sensor's reporting mode as returned by ASensor_getReportingMode.
type ReportingMode int32

const (
	ReportingInvalid        ReportingMode = -1
	ReportingContinuous     ReportingMode = 0
	ReportingOnChange       ReportingMode = 1
	ReportingOneShot        ReportingMode = 2
	ReportingSpecialTrigger ReportingMode = 3
)

// String returns the Android name of the reporting mode.
func (m ReportingMode) String() string {
	switch m {
	case ReportingContinuous:
		return "continuous"
	case ReportingOnChange:
		return "on-change"
	case ReportingOneShot:
		return "one-shot"
	case ReportingSpecialTrigger:
		return "special-trigger"
	default:
		return fmt.Sprintf("ReportingMode(%d)", int(m))
	}
}

// Manager wraps a native ASensorManager handle.
type Manager struct {
	ptr uintptr
}

// GetInstance returns the process-wide sensor manager.
func GetInstance() (*Manager, error) {
	if err := load(); err != nil {
		return nil, err
	}
	ptr := fnGetInstance()
	if ptr == 0 {
		return nil, fmt.Errorf("ASensorManager_getInstance returned NULL")
	}
	return &Manager{ptr: ptr}, nil
}

// Sensors returns every sensor reported by the manager.
func (m *Manager) Sensors() ([]*Sensor, error) {
	var list unsafe.Pointer
	count := fnGetSensorList(m.ptr, &list)
	if count < 0 {
		return nil, fmt.Errorf("ASensorManager_getSensorList returned %d", count)
	}

	const ptrSize = unsafe.Sizeof(uintptr(0))
	sensors := make([]*Sensor, 0, count)
	for i := int32(0); i < count; i++ {
		ref := *(*uintptr)(unsafe.Add(list, uintptr(i)*ptrSize))
		sensors = append(sensors, &Sensor{ptr: ref})
	}
	return sensors, nil
}

// startQueue pins the goroutine to its OS thread (the Android looper is
// thread-local), prepares a looper, creates an event queue, and enables s on
// it. The returned cleanup disables the sensor, releases the queue and
// looper, and unpins the thread.
func (m *Manager) startQueue(s *Sensor) (uintptr, func(), error) {
	runtime.LockOSThread()

	looper := fnALooperPrepare(1) // ALOOPER_PREPARE_ALLOW_NON_CALLBACKS
	if looper == 0 {
		runtime.UnlockOSThread()
		return 0, nil, fmt.Errorf("ALooper_prepare returned NULL")
	}

	queue := fnCreateEventQueue(m.ptr, looper, 0, 0, 0)
	if queue == 0 {
		fnALooperRelease(looper)
		runtime.UnlockOSThread()
		return 0, nil, fmt.Errorf("ASensorManager_createEventQueue returned NULL")
	}

	if rc := fnEnableSensor(queue, s.ptr); rc != 0 {
		fnDestroyEventQueue(m.ptr, queue)
		fnALooperRelease(looper)
		runtime.UnlockOSThread()
		return 0, nil, fmt.Errorf("enable %s: ASensorEventQueue_enableSensor returned %d", s.Name(), rc)
	}

	cleanup := func() {
		fnDisableSensor(queue, s.ptr)
		fnDestroyEventQueue(m.ptr, queue)
		fnALooperRelease(looper)
		runtime.UnlockOSThread()
	}
	return queue, cleanup, nil
}

// pollEvent blocks on the looper until deadline for a single event. A zero
// deadline means wait forever. ok is false when the deadline passes without
// an event.
func pollEvent(queue uintptr, deadline time.Time) (*Event, bool, error) {
	var raw cEvent
	for {
		rc := fnHasEvents(queue)
		if rc == 1 {
			n := fnGetEvents(queue, &raw, 1)
			if n == 1 {
				return &Event{Sensor: raw.Sensor, Type: Type(raw.Type), Timestamp: raw.Timestamp, Data: raw.Data}, true, nil
			}
			if n < 0 {
				return nil, false, fmt.Errorf("ASensorEventQueue_getEvents returned %d", n)
			}
			// n == 0: readiness was spurious, keep waiting.
		} else if rc < -1 && rc != -4 {
			// The NDK implementation returns -1 (not 0) when no event is
			// pending, and a value < -1 for a real error. -4 is EINTR,
			// which the Go runtime can trigger; treat it as transient.
			return nil, false, fmt.Errorf("ASensorEventQueue_hasEvents returned %d", rc)
		}

		// Block until an event is ready or the deadline passes.
		timeoutMillis := int32(-1) // wait indefinitely
		if !deadline.IsZero() {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return nil, false, nil
			}
			// Round up so the poll never wakes before the deadline.
			timeoutMillis = int32((remaining + time.Millisecond - 1) / time.Millisecond)
		}

		if rc := fnALooperPollOnce(timeoutMillis, 0, 0, 0); rc == -4 {
			// ALOOPER_POLL_ERROR: the calling thread has no looper or an
			// unrecoverable error occurred.
			return nil, false, fmt.Errorf("ALooper_pollOnce returned %d", rc)
		}
		// -1 (wake), -2 (callback), -3 (timeout), or an ident (event ready)
		// all loop around to drain the queue and re-check the deadline.
	}
}

// Read enables the sensor, waits for a single event, and disables it again.
func (m *Manager) Read(s *Sensor, timeout time.Duration) (*Event, error) {
	queue, cleanup, err := m.startQueue(s)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	ev, ok, err := pollEvent(queue, time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("no event within %s", timeout)
	}
	return ev, nil
}

// Watch enables the sensor and streams events to fn until duration elapses.
// A non-positive duration streams indefinitely.
func (m *Manager) Watch(s *Sensor, duration time.Duration, fn func(*Event)) error {
	queue, cleanup, err := m.startQueue(s)
	if err != nil {
		return err
	}
	defer cleanup()

	var deadline time.Time
	if duration > 0 {
		deadline = time.Now().Add(duration)
	}
	for {
		ev, ok, err := pollEvent(queue, deadline)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		fn(ev)
	}
}

// Sensor wraps a native ASensor handle.
type Sensor struct {
	ptr uintptr
}

// Name returns the sensor's name.
func (s *Sensor) Name() string { return fnGetName(s.ptr) }

// Vendor returns the sensor's vendor string.
func (s *Sensor) Vendor() string { return fnGetVendor(s.ptr) }

// Type returns the sensor's type.
func (s *Sensor) Type() Type { return Type(fnGetType(s.ptr)) }

// StringType returns the sensor's Android string type (e.g. android.sensor.accelerometer).
func (s *Sensor) StringType() string { return fnGetStringType(s.ptr) }

// MinDelay returns the minimum delay between events in microseconds.
func (s *Sensor) MinDelay() int32 { return fnGetMinDelay(s.ptr) }

// Resolution returns the sensor's resolution in the sensor's unit.
func (s *Sensor) Resolution() float32 { return fnGetResolution(s.ptr) }

// ReportingMode returns the sensor's reporting mode.
func (s *Sensor) ReportingMode() ReportingMode { return ReportingMode(fnReportingMode(s.ptr)) }

// FifoMaxEventCount returns the maximum number of events that can be batched.
func (s *Sensor) FifoMaxEventCount() int32 { return fnFifoMax(s.ptr) }

// FifoReservedEventCount returns the number of events reserved in the FIFO.
func (s *Sensor) FifoReservedEventCount() int32 { return fnFifoReserved(s.ptr) }

// Handle returns the sensor's system handle.
func (s *Sensor) Handle() int32 { return fnGetHandle(s.ptr) }

// HighestDirectReportRateLevel returns the highest direct report rate level supported.
func (s *Sensor) HighestDirectReportRateLevel() int32 { return fnDirectRate(s.ptr) }

// IsWakeUp reports whether this is a wake-up sensor.
func (s *Sensor) IsWakeUp() bool { return fnIsWakeUpSensor(s.ptr) }

// Event is a single sensor reading.
type Event struct {
	Sensor    int32
	Type      Type
	Timestamp int64 // elapsed realtime nanoseconds
	Data      [16]float32
}

// StepCount returns the step counter value for a step counter event.
func (e *Event) StepCount() uint64 {
	return *(*uint64)(unsafe.Pointer(&e.Data[0]))
}
