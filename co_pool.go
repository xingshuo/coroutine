package coroutine

type CoNode struct {
	co   *Coroutine
	next *CoNode
}

// 栈式单链表
var freeNodeList = &CoNode{}

func pushToNodeList(co *Coroutine) {
	freeNodeList.next = &CoNode{
		co:   co,
		next: freeNodeList.next,
	}
}

func popFromNodeList() *Coroutine {
	var co *Coroutine = nil
	next := freeNodeList.next
	if next != nil {
		co = next.co
		freeNodeList.next = next.next
		// 这里不处理也不影响gc, 严谨行事
		next.next = nil
	}
	return co
}

func sizeOfNodeList() int {
	cur := freeNodeList.next
	sz := 0
	for cur != nil {
		cur = cur.next
		sz++
	}
	return sz
}

func CreateFromPool(f CoFunc) *Coroutine {
	nco := popFromNodeList()
	if nco == nil {
		// 若执行f过程中发生错误, 则coroutine永远不会pushToFreeList, 交给golang自动gc清理
		nco = Create(func(co *Coroutine, a ...interface{}) []interface{} {
			f(co, a...)
			for {
				f = nil
				pushToNodeList(co)
				f = co.Yield()[0].(CoFunc)
				f(co, co.Yield()...)
			}
		})
	} else {
		nco.Resume(f)
	}
	return nco
}
