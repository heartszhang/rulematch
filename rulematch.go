package rulematch

import (
	ac "github.com/cloudflare/ahocorasick"
)

type bits []byte
type keyword_rule bits

type keyword_matcher struct {
	matcher   *ac.Matcher
	_words    map[string]struct{}
	words     []string
	rules     []keyword_rule
	word2rule [][]int
}

type Matcher interface {
	Match(line string) []int // rule list
}

func NewMatcher(rules ...[]string) Matcher {
	this := &keyword_matcher{_words: map[string]struct{}{}}
	for idx, rule := range rules {
		this.rules = append(this.rules, keyword_rule{})
		for _, word := range rule {
			if _, ok := this._words[word]; !ok {
				this._words[word] = struct{}{}
				this.words = append(this.words, word)
				(*bits)(&this.rules[idx]).set(len(this.word2rule))
				this.word2rule = append(this.word2rule, []int{idx})
			}
		}
	}
	this.matcher = ac.NewStringMatcher(this.words)
	return this
}

func (this *keyword_matcher) Match(line string) (v []int) {
	words := this.matcher.Match([]byte(line))
	var whole, r = word2bits(words), bits{}
	for _, w := range words {
		for _, rule := range this.word2rule[w] {
			if !r.has(rule) {
				r.set(rule)
				b := this.rules[rule]
				if whole.intersect(bits(b)) == string(b) {
					v = append(v, rule)
				}
			}
		}
	}
	return
}

/*
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
*/
