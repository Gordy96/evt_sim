package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParameters(t *testing.T) {
	var p Parameters

	p.WithBool("bool", true)
	p.WithInt64("int64", int64(100))
	p.WithFloat64("float64", 1.23)
	p.WithString("string", "hello")
	p.WithDuration("duration", time.Second)
	now := time.Now()
	p.WithTime("time", now)
	p.WithStrings("string_list", []string{"hello", "world"})
	p.SetInt32("int32", int32(100))

	i32, ok := p.GetInt32("int32")
	assert.True(t, ok)
	assert.Equal(t, int32(100), i32)

	assert.Equal(t, map[string]any{
		"bool":        true,
		"int64":       int64(100),
		"float64":     1.23,
		"string":      "hello",
		"duration":    time.Second,
		"time":        now,
		"string_list": []string{"hello", "world"},
		"int32":       int32(100),
	}, p.params)
}

func TestParameters_Merge(t *testing.T) {
	var p Parameters
	p.SetString("string", "hello")

	var c Parameters

	c.SetStrings("string_list", []string{"hello", "world"})

	c.Merge(p)

	assert.Equal(t, map[string]any{
		"string":      "hello",
		"string_list": []string{"hello", "world"},
	}, c.params)
}
