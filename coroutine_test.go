package coroutine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	local count = 0
	local co = coroutine.create(function (a, b)
		while count < 2 do
			local c, d = a + b, a - b
			a, b = coroutine.yield(c, d)
			count = count + 1
		end
		assert(a == "aaa" and b == "bbb")
		local e, f = coroutine.yield("wasd", 666)
		assert(e == 700 and f == "tencent")
		-- here occurred error
		return e / f
	end)

	local ok1, sum1, sub1 = coroutine.resume(co, 100, 10)
	assert(ok1 == true and sum1 == 110 and sub1 == 90)

	local ok2, sum2, sub2 = coroutine.resume(co, sum1, sub1)
	assert(ok2 == true and sum2 == 200 and sub2 == 20)

	local ok3, r1, r2 = coroutine.resume(co, "aaa", "bbb")
	assert(ok3 == true and r1 == "wasd" and r2 == 666)

	local ok4, err = coroutine.resume(co, 700, "tencent")
	assert(ok4 == false and coroutine.status(co) == "dead")
*/

// 该测试逻辑与上面注释中的lua代码基本保持一致
func TestCoroutine(t *testing.T) {
	count := 0
	co := Create(func(co *Coroutine, args ...interface{}) []interface{} {
		for count < 2 {
			a, b := args[0].(int), args[1].(int)
			c, d := a+b, a-b
			args = co.Yield(c, d)
			count++
		}
		as, bs := args[0].(string), args[1].(string)
		assert.Equal(t, "aaa", as)
		assert.Equal(t, "bbb", bs)

		args = co.Yield("wasd", 666)
		e, f := args[0].(int), args[1].(int)
		assert.Equal(t, 700, e)
		assert.Equal(t, 0, f)
		// here occurred error
		return []interface{}{e / f}
	})

	err, retvals := co.Resume(100, 10)
	assert.Nil(t, err)
	assert.Equal(t, 110, retvals[0])
	assert.Equal(t, 90, retvals[1])

	err, retvals = co.Resume(retvals[0], retvals[1])
	assert.Nil(t, err)
	assert.Equal(t, 200, retvals[0])
	assert.Equal(t, 20, retvals[1])

	err, retvals = co.Resume("aaa", "bbb")
	assert.Nil(t, err)
	assert.Equal(t, "wasd", retvals[0])
	assert.Equal(t, 666, retvals[1])

	err, retvals = co.Resume(700, 0)
	assert.NotNil(t, err)
	assert.Equal(t, Dead, co.Status())
	assert.Nil(t, retvals)
}
