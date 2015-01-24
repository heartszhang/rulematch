package main

import (
	"fmt"

	ac "github.com/cloudflare/ahocorasick"
	rm "github.com/heartszhang/rulematch"
)

var sample = `60.165.25.81 - - [20/Jan/2015:13:00:03 +0800] "GET /firstGame/Android/funs0004/0.1.0.60_0.0.0.0.1/serverUnknown/3ebb9cc16e4d5ea6d3ffb4abe6b51d668/userUnknown/1/0/Start HTTP/1.1" 200 151 "-" "-"`

func main() {
	/*	km := keyword_matcher{_words: map[string]struct{}{}}
		km.add_rule(nil, "HTTP", "get")
		km.add_rule(nil, "HTTP", "GET")
		km.add_rule(nil, "GET")
		km.add_rule(nil, "WHAT", "none")
		km.add_rule(nil, "GET", "firstGame")
		km.build()
		rules := km.match(sample)
	*/
	km := rm.NewMatcher([]string{"HTTP", "get"},
		[]string{"HTTP", "GET"}, []string{"GET"},
		[]string{"WHAT", "none"}, []string{"GET", "firstGame"})
	rules := km.Match(sample)
	fmt.Println(rules)
}

type rule interface {
	decode(content string) packet
}

type packet map[string]interface{}
type bits []byte

type keyword_rule struct {
	decoder func(string) packet
	bits    bits
}

type matcher interface {
	add_rule(decoder func(string) packet, words ...string)
	match(line string) []int
}

type keyword_matcher struct {
	matcher   *ac.Matcher
	_words    map[string]struct{}
	words     []string
	rules     []keyword_rule
	word2rule [][]int
}

func (this *keyword_matcher) add_rule(decoder func(string) packet, words ...string) {
	idx := len(this.rules)
	this.rules = append(this.rules, keyword_rule{decoder: decoder})
	for _, word := range words {
		if _, ok := this._words[word]; !ok {
			this._words[word] = struct{}{}
			this.words = append(this.words, word)
			this.rules[idx].bits.set(len(this.word2rule))
			this.word2rule = append(this.word2rule, []int{idx})
		}
	}
}
func (this *keyword_matcher) build() {
	fmt.Println(this.words)
	fmt.Println(this._words)
	this.matcher = ac.NewStringMatcher(this.words)
}

func (this *keyword_matcher) match(line string) (v []int) {
	words := this.matcher.Match([]byte(line))
	fmt.Println("keywords", words)
	var whole, r = word2bits(words), bits{}
	for _, w := range words {
		for _, rule := range this.word2rule[w] {
			if !r.has(rule) {
				r.set(rule)
				b := this.rules[rule].bits
				if whole.intersect(b) == string(b) {
					v = append(v, rule)
				}
			}
		}
	}
	return
}

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
