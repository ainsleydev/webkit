package ptr

import "time"

// String returns the string value of a pointer to string.
// If the pointer is nil, it returns an empty string.
func String(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// StringPtr returns a pointer to the provided string value.
func StringPtr(v string) *string {
	return &v
}

// Byte returns the byte value of a pointer to byte.
// If the pointer is nil, it returns 0.
func Byte(v *byte) byte {
	if v == nil {
		return 0
	}
	return *v
}

// BytePtr returns a pointer to the provided byte value.
func BytePtr(v byte) *byte {
	return &v
}

// Bool returns the bool value of a pointer to bool.
// If the pointer is nil, it returns false.
func Bool(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

// BoolPtr returns a pointer to the provided bool value.
func BoolPtr(v bool) *bool {
	return &v
}

// Int returns the int value of a pointer to int.
// If the pointer is nil, it returns 0.
func Int(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// IntPtr returns a pointer to the provided int value.
func IntPtr(v int) *int {
	return &v
}

// Int8 returns the int8 value of a pointer to int8.
// If the pointer is nil, it returns 0.
func Int8(v *int8) int8 {
	if v == nil {
		return 0
	}
	return *v
}

// Int8Ptr returns a pointer to the provided int8 value.
func Int8Ptr(v int8) *int8 {
	return &v
}

// Int16 returns the int16 value of a pointer to int16.
// If the pointer is nil, it returns 0.
func Int16(v *int16) int16 {
	if v == nil {
		return 0
	}
	return *v
}

// Int16Ptr returns a pointer to the provided int16 value.
func Int16Ptr(v int16) *int16 {
	return &v
}

// Int32 returns the int32 value of a pointer to int32.
// If the pointer is nil, it returns 0.
func Int32(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

// Int32Ptr returns a pointer to the provided int32 value.
func Int32Ptr(v int32) *int32 {
	return &v
}

// Int64 returns the int64 value of a pointer to int64.
// If the pointer is nil, it returns 0.
func Int64(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

// Int64Ptr returns a pointer to the provided int64 value.
func Int64Ptr(v int64) *int64 {
	return &v
}

// Float32 returns the float32 value of a pointer float32.
// If the pointer is nil, it returns 0.
func Float32(v *float32) float32 {
	if v == nil {
		return 0
	}
	return *v
}

// Float32Ptr returns a pointer to the provided float32 value.
func Float32Ptr(v float32) *float32 {
	return &v
}

// Float64 returns the float64 value of a pointer float64.
// If the pointer is nil, it returns 0.
func Float64(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

// Float64Ptr returns a pointer to the provided float64 value.
func Float64Ptr(v float64) *float64 {
	return &v
}

// Uint returns the uint value of a pointer to uint.
// If the pointer is nil, it returns 0.
func Uint(v *uint) uint {
	if v == nil {
		return 0
	}
	return *v
}

// UintPtr returns a pointer to the provided uint value.
func UintPtr(v uint) *uint {
	return &v
}

// Uint8 returns the uint8 value of a pointer to uint8.
// If the pointer is nil, it returns 0.
func Uint8(v *uint8) uint8 {
	if v == nil {
		return 0
	}
	return *v
}

// Uint8Ptr returns a pointer to the provided uint8 value.
func Uint8Ptr(v uint8) *uint8 {
	return &v
}

// Uint16 returns the uint16 value of a pointer to uint16.
// If the pointer is nil, it returns 0.
func Uint16(v *uint16) uint16 {
	if v == nil {
		return 0
	}
	return *v
}

// Uint16Ptr returns a pointer to the provided uint16 value.
func Uint16Ptr(v uint16) *uint16 {
	return &v
}

// Uint32 returns the uint32 value of a pointer to uint32.
// If the pointer is nil, it returns 0.
func Uint32(v *uint32) uint32 {
	if v == nil {
		return 0
	}
	return *v
}

// Uint32Ptr returns a pointer to the provided uint32 value.
func Uint32Ptr(v uint32) *uint32 {
	return &v
}

// Uint64 returns the uint64 value of a pointer to uint64.
// If the pointer is nil, it returns 0.
func Uint64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

// Uint64Ptr returns a pointer to the provided uint64 value.
func Uint64Ptr(v uint64) *uint64 {
	return &v
}

// Uintptr returns the uintptr value of a pointer to uintptr.
// If the pointer is nil, it returns 0.
func Uintptr(v *uintptr) uintptr {
	if v == nil {
		return 0
	}
	return *v
}

// UintptrPtr returns a pointer to the provided uintptr value.
func UintptrPtr(v uintptr) *uintptr {
	return &v
}

// Complex64 returns the complex64 value of a pointer to complex64.
// If the pointer is nil, it returns 0+0i.
func Complex64(v *complex64) complex64 {
	if v == nil {
		return 0 + 0i
	}
	return *v
}

// Complex64Ptr returns a pointer to the provided complex64 value.
func Complex64Ptr(v complex64) *complex64 {
	return &v
}

// Complex128 returns the complex128 value of a pointer to complex128.
// If the pointer is nil, it returns 0+0i.
func Complex128(v *complex128) complex128 {
	if v == nil {
		return 0 + 0i
	}
	return *v
}

// Complex128Ptr returns a pointer to the provided complex128 value.
func Complex128Ptr(v complex128) *complex128 {
	return &v
}

// Time returns the time.RecordedAt value of a pointer to time.RecordedAt.
// If the pointer is nil, it returns the zero value of time.RecordedAt.
func Time(v *time.Time) time.Time {
	if v == nil {
		return time.Time{}
	}
	return *v
}

// TimePtr returns a pointer to the provided time.RecordedAt value.
func TimePtr(v time.Time) *time.Time {
	return &v
}
