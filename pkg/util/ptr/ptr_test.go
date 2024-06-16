package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var str *string
		assert.Equal(t, "", String(str))
	})

	t.Run("Non Nil", func(t *testing.T) {
		s := "test"
		str := &s
		assert.Equal(t, s, String(str))
	})
}

func TestStringPtr(t *testing.T) {
	s := "test"
	ptr := StringPtr(s)
	assert.Equal(t, s, *ptr)
}

func TestByte(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var b *byte
		assert.Equal(t, byte(0), Byte(b))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := byte(10)
		ptr := &val
		assert.Equal(t, val, Byte(ptr))
	})
}

func TestBytePtr(t *testing.T) {
	val := byte(10)
	ptr := BytePtr(val)
	assert.Equal(t, val, *ptr)
}

func TestBool(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var b *bool
		assert.False(t, Bool(b))
	})

	t.Run("True", func(t *testing.T) {
		b := true
		ptr := &b
		assert.True(t, Bool(ptr))
	})

	t.Run("False", func(t *testing.T) {
		b := false
		ptr := &b
		assert.False(t, Bool(ptr))
	})
}

func TestBoolPtr(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		val := true
		ptr := BoolPtr(val)
		assert.True(t, *ptr)
	})

	t.Run("False", func(t *testing.T) {
		val := false
		ptr := BoolPtr(val)
		assert.False(t, *ptr)
	})
}

func TestInt(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var i *int
		assert.Equal(t, 0, Int(i))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := 42
		ptr := &val
		assert.Equal(t, val, Int(ptr))
	})
}

func TestIntPtr(t *testing.T) {
	val := 42
	ptr := IntPtr(val)
	assert.Equal(t, val, *ptr)
}

func TestInt8(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var i *int8
		assert.Equal(t, int8(0), Int8(i))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := int8(42)
		ptr := &val
		assert.Equal(t, val, Int8(ptr))
	})
}

func TestInt8Ptr(t *testing.T) {
	val := int8(42)
	ptr := Int8Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestInt16(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var i *int16
		assert.Equal(t, int16(0), Int16(i))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := int16(42)
		ptr := &val
		assert.Equal(t, val, Int16(ptr))
	})
}

func TestInt16Ptr(t *testing.T) {
	val := int16(42)
	ptr := Int16Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestInt32(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var i *int32
		assert.Equal(t, int32(0), Int32(i))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := int32(42)
		ptr := &val
		assert.Equal(t, val, Int32(ptr))
	})
}

func TestInt32Ptr(t *testing.T) {
	val := int32(42)
	ptr := Int32Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestInt64(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var i *int64
		assert.Equal(t, int64(0), Int64(i))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := int64(42)
		ptr := &val
		assert.Equal(t, val, Int64(ptr))
	})
}

func TestInt64Ptr(t *testing.T) {
	val := int64(42)
	ptr := Int64Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestFloat32(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var i *float32
		assert.Equal(t, float32(0), Float32(i))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := float32(42.2)
		ptr := &val
		assert.Equal(t, val, Float32(ptr))
	})
}

func TestFloat32Ptr(t *testing.T) {
	val := float32(42.2)
	ptr := Float32Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestFloat64(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var i *float64
		assert.Equal(t, float64(0), Float64(i))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := 42.2
		ptr := &val
		assert.Equal(t, val, Float64(ptr))
	})
}

func TestFloat64Ptr(t *testing.T) {
	val := 42.2
	ptr := Float64Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestUint(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var u *uint
		assert.Equal(t, uint(0), Uint(u))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := uint(42)
		ptr := &val
		assert.Equal(t, val, Uint(ptr))
	})
}

func TestUintPtr(t *testing.T) {
	val := uint(42)
	ptr := UintPtr(val)
	assert.Equal(t, val, *ptr)
}

func TestUint8(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var u *uint8
		assert.Equal(t, uint8(0), Uint8(u))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := uint8(42)
		ptr := &val
		assert.Equal(t, val, Uint8(ptr))
	})
}

func TestUint8Ptr(t *testing.T) {
	val := uint8(42)
	ptr := Uint8Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestUint16(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var u *uint16
		assert.Equal(t, uint16(0), Uint16(u))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := uint16(42)
		ptr := &val
		assert.Equal(t, val, Uint16(ptr))
	})
}

func TestUint16Ptr(t *testing.T) {
	val := uint16(42)
	ptr := Uint16Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestUint32(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var u *uint32
		assert.Equal(t, uint32(0), Uint32(u))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := uint32(42)
		ptr := &val
		assert.Equal(t, val, Uint32(ptr))
	})
}

func TestUint32Ptr(t *testing.T) {
	val := uint32(42)
	ptr := Uint32Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestUint64(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var u *uint64
		assert.Equal(t, uint64(0), Uint64(u))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := uint64(42)
		ptr := &val
		assert.Equal(t, val, Uint64(ptr))
	})
}

func TestUint64Ptr(t *testing.T) {
	val := uint64(42)
	ptr := Uint64Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestUintptr(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var u *uintptr
		assert.Equal(t, uintptr(0), Uintptr(u))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := uintptr(42)
		ptr := &val
		assert.Equal(t, val, Uintptr(ptr))
	})
}

func TestUintptrPtr(t *testing.T) {
	val := uintptr(42)
	ptr := UintptrPtr(val)
	assert.Equal(t, val, *ptr)
}

func TestComplex64(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var c *complex64
		assert.Equal(t, complex64(0+0i), Complex64(c))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := complex64(3 + 4i)
		ptr := &val
		assert.Equal(t, val, Complex64(ptr))
	})
}

func TestComplex64Ptr(t *testing.T) {
	val := complex64(3 + 4i)
	ptr := Complex64Ptr(val)
	assert.Equal(t, val, *ptr)
}

func TestComplex128(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var c *complex128
		assert.Equal(t, complex128(0+0i), Complex128(c))
	})

	t.Run("Non Nil", func(t *testing.T) {
		val := complex128(3 + 4i)
		ptr := &val
		assert.Equal(t, val, Complex128(ptr))
	})
}

func TestComplex128Ptr(t *testing.T) {
	val := complex128(3 + 4i)
	ptr := Complex128Ptr(val)
	assert.Equal(t, val, *ptr)
}
