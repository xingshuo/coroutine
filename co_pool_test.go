package coroutine

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 作者正在经历失眠...
func CountingSheep(co *Coroutine, a ...interface{}) []interface{} {
	user, nSheep := a[0].(string), a[1].(int)
	cnt := 0
	for cnt < nSheep {
		co.Yield()
		cnt++
		log.Printf("%s counting sheep %d times", user, cnt)
	}
	log.Printf("%s fall asleep!", user)
	return nil
}

func TestCoroutinePool(t *testing.T) {
	coA := CreateFromPool(CountingSheep)
	coA.Resume("workerA", 1)

	coB := CreateFromPool(CountingSheep)
	coB.Resume("workerB", 2)

	coC := CreateFromPool(CountingSheep)
	coC.Resume("workerC", 1)
	assert.Equal(t, 0, sizeOfNodeList())

	coA.Resume()
	assert.Equal(t, 1, sizeOfNodeList())

	coD := CreateFromPool(CountingSheep)
	assert.Equal(t, coA, coD)

	coD.Resume("workerD", 2)
	assert.Equal(t, 0, sizeOfNodeList())

	coB.Resume()
	assert.Equal(t, 0, sizeOfNodeList())

	coB.Resume()
	assert.Equal(t, 1, sizeOfNodeList())

	coC.Resume()
	assert.Equal(t, 2, sizeOfNodeList())

	coD.Resume()
	assert.Equal(t, 2, sizeOfNodeList())

	coD.Resume()
	assert.Equal(t, 3, sizeOfNodeList())

	assert.Equal(t, coD, freeNodeList.next.co)
	assert.Equal(t, coC, freeNodeList.next.next.co)
	assert.Equal(t, coB, freeNodeList.next.next.next.co)
	assert.Equal(t, (*CoNode)(nil), freeNodeList.next.next.next.next)
}
