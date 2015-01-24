package rulematch

func (lhs bits) intersect(rhs bits) string {
	if len(lhs) < len(rhs) {
		return ""
	}
	v := make([]byte, len(rhs))
	for idx, val := range rhs {
		v[idx] = val & lhs[idx]
	}
	return string(v)
}

func word2bits(words []int) (v bits) {
	for _, idx := range words {
		v.set(idx)
	}
	return
}

func (this *bits) set(idx int) {
	if len(*this)*8 <= idx {
		n := make([]byte, idx/8+1)
		copy(n, *this)
		*this = n
	}
	x, off := idx/8, idx%8
	(*this)[x] |= 1 << uint(off)
}

func (this bits) has(idx int) bool {
	x, off := idx/8, idx%8
	if x >= len(this) {
		return false
	}
	return this[x]&(1<<uint(off)) == 1<<uint(off)
}
