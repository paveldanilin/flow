package flow

type ObjectPool struct {
	values chan interface{}
	alloc  func() interface{}
}

func NewObjectPool(size int) *ObjectPool {
	return &ObjectPool{
		values: make(chan interface{}, size),
	}
}

func (p *ObjectPool) Init(alloc func() interface{}) {
	p.alloc = alloc

	for i := 0; i < cap(p.values); i++ {
		p.values <- p.alloc()
	}
}

func (p *ObjectPool) Get() (interface{}, bool) {
	var v interface{}
	select {
	case v = <-p.values:
		return v, true
	default:
		// We ran out of the pool capacity, just allocate a new object. In this case we don't need to put it back.
		v = p.alloc()
		return v, false
	}
}

func (p *ObjectPool) Put(obj interface{}) {
	select {
	case p.values <- obj:
	default:
	}
}
