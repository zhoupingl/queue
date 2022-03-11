package queue

// 最新二叉堆

type BinaryHead interface {
	// 添加一个原型
	Push(class int, val interface{})
	// 弹出一个数据
	Pop() (*Item, bool)
}

func NewBinaryHead() BinaryHead {
	l := &BinaryHeadImpl{}
	l.Init()

	return l
}

type Item struct {
	Class int
	Val   interface{}
}

type BinaryHeadImpl struct {
	list map[int]*Item
}

func (b *BinaryHeadImpl) Init() {
	b.list = make(map[int]*Item, 1000)
}

func (b *BinaryHeadImpl) Len() int {
	return len(b.list)
}

// 从1,2,3...n
func (b *BinaryHeadImpl) NextAddr() int {
	return b.Len() + 1
}

func (b *BinaryHeadImpl) ParentAddr(addr int) int {
	return addr / 2
}

// 插入元素时候使用
func (b *BinaryHeadImpl) ShiftUp(addr int) {
	parentAddr := b.ParentAddr(addr)
	if !b.ValidAddr(parentAddr) {
		return
	}

	if b.compareClass(parentAddr, addr) == 1 {
		b.Swap(parentAddr, addr)
		b.ShiftUp(parentAddr)
	}

}

// 比较两个节点优先级
// addr < addr2 return -1
// addr = addr2 return 0
// addr > addr2 return 1
func (b *BinaryHeadImpl) compareClass(addr, addr2 int) int {
	if b.list[addr].Class < b.list[addr2].Class {
		return -1
	} else if b.list[addr].Class == b.list[addr2].Class {
		return 0
	} else {
		return 1
	}
}

// 交互位置
func (b *BinaryHeadImpl) Swap(addr, addr2 int) {
	b.list[addr], b.list[addr2] = b.list[addr2], b.list[addr]
}

// 求左右节点 最小节点位置
func (b *BinaryHeadImpl) LeafMinAddr(addr int) int {
	leftAddr := b.LeftAddr(addr)
	rightAddr := b.RightAddr(addr)

	if b.ValidAddr(leftAddr) && b.ValidAddr(rightAddr) { // 情况1 存在左右节点
		if b.compareClass(leftAddr, rightAddr) == -1 {
			return leftAddr
		}
		return rightAddr
	} else if b.ValidAddr(leftAddr) { // 情况2 只有左节点
		return leftAddr
	} else { // 情况3 没有左右节点
		return 0
	}
}

// 弹出一个元素时候使用
func (b *BinaryHeadImpl) ShiftDown(addr int) {

	swapAddr := b.LeafMinAddr(addr)
	if !b.ValidAddr(swapAddr) { // 不存在交换位置
		return
	}

	if b.compareClass(addr, swapAddr) == 1 {
		b.Swap(addr, swapAddr)
		b.ShiftDown(swapAddr)
	}
}

// 判定是否有效节点位置
func (b *BinaryHeadImpl) ValidAddr(addr int) bool {
	if addr == 0 || addr > b.Len() { // 无效节点位置
		return false
	}

	return true
}

func (b *BinaryHeadImpl) LeftAddr(addr int) int {
	return addr * 2
}

func (b *BinaryHeadImpl) RightAddr(addr int) int {
	return b.LeftAddr(addr) + 1
}

func (b *BinaryHeadImpl) Push(class int, val interface{}) {

	// 添加一个元素
	node := &Item{
		Class: class,
		Val:   val,
	}
	b.list[b.NextAddr()] = node

	// 调用shift up
	b.ShiftUp(b.Len())
}

func (b *BinaryHeadImpl) Pop() (*Item, bool) {

	if b.Len() == 0 { // 没有数据
		return nil, false
	}

	// 获取top节点
	topNode := b.list[1]
	// 最后一个节点，移动到top节点
	b.list[1] = b.list[b.Len()]
	// 销毁最后元素位置
	delete(b.list, b.Len())

	// 调用shit down
	if b.ValidAddr(1) {
		b.ShiftDown(1)
	}

	return topNode, true
}
