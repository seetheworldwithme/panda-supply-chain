package warehouse

// HanLP portable-1.8.4 (Apache-2.0) 默认 ViterbiSegment 的原生移植。
// 模型由同版本原 JAR 离线导出；词网、双词权重、自定义词典及人名识别均保留原顺序。
// 来源、许可与可复现命令见 hanlp_data/README.md；生产运行不需要 JVM。
import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"sync"
	"unicode/utf16"
)

//go:embed hanlp_data/model.bin.gz
var hanlpModelData []byte

type hlLabel struct{ tag, frequency int }
type hlAttr struct {
	total  int
	labels []hlLabel
}
type hlDAT struct {
	base, check []int32
	attrs       []*hlAttr
	roles       [][]hlLabel
}
type hlModel struct {
	total, oov            int
	myu                   float64
	normalization         bool
	natures               []string
	charType              []byte
	charTable             []uint16
	core, custom, nr, nrf hlDAT
	start, pair           []int32
	nrTotal               [16]int
	nrTransition          [16][16]float64
	ids                   map[string]int
}
type hlVertex struct {
	real   []uint16
	attr   *hlAttr
	id     int
	from   *hlVertex
	weight float64
}
type hlNet [][]*hlVertex

var loadHanlp = sync.OnceValue(func() *hlModel {
	z, err := gzip.NewReader(bytes.NewReader(hanlpModelData))
	if err != nil {
		panic(err)
	}
	defer z.Close()
	data, err := io.ReadAll(z)
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(data)
	read := func(v any) {
		if err := binary.Read(r, binary.BigEndian, v); err != nil {
			panic(fmt.Errorf("HanLP 模型损坏: %w", err))
		}
	}
	i32 := func() int { var n int32; read(&n); return int(n) }
	text := func() string {
		n := i32()
		b := make([]byte, n)
		if _, err = io.ReadFull(r, b); err != nil {
			panic(err)
		}
		return string(b)
	}
	header := make([]byte, 6)
	io.ReadFull(r, header)
	if string(header) != "HLP184" {
		panic("HanLP 模型版本不匹配")
	}
	m := &hlModel{total: i32(), ids: map[string]int{}}
	read(&m.myu)
	m.oov = i32()
	read(&m.normalization)
	m.natures = make([]string, i32())
	for i := range m.natures {
		m.natures[i] = text()
	}
	m.charType = make([]byte, 65536)
	read(m.charType)
	m.charTable = make([]uint16, 65536)
	read(m.charTable)
	dat := func(d *hlDAT) {
		n := i32()
		d.base = make([]int32, n+65536)
		d.check = make([]int32, n+65536)
		for i := 0; i < n; i++ {
			read(&d.base[i])
			read(&d.check[i])
		}
	}
	attrs := func(d *hlDAT) {
		d.attrs = make([]*hlAttr, i32())
		for i := range d.attrs {
			a := &hlAttr{total: i32()}
			a.labels = make([]hlLabel, i32())
			for j := range a.labels {
				a.labels[j] = hlLabel{i32(), i32()}
			}
			d.attrs[i] = a
		}
		dat(d)
	}
	attrs(&m.core)
	attrs(&m.custom)
	ints := func() []int32 { a := make([]int32, i32()); read(a); return a }
	m.start = ints()
	m.pair = ints()
	m.nr.roles = make([][]hlLabel, i32())
	for i := range m.nr.roles {
		m.nr.roles[i] = make([]hlLabel, i32())
		for j := range m.nr.roles[i] {
			m.nr.roles[i][j] = hlLabel{i32(), i32()}
		}
	}
	dat(&m.nr)
	for i := range m.nrTotal {
		m.nrTotal[i] = i32()
	}
	read(&m.nrTransition)
	dat(&m.nrf)
	if r.Len() != 0 {
		panic("HanLP 模型存在未读取数据")
	}
	for _, word := range []string{"始##始", "末##末", "未##人", "未##地", "未##团", "未##数", "未##串", "未##时", "未##专"} {
		m.ids[word] = m.core.lookup(hlChars(word))
	}
	return m
})

func hlChars(s string) []uint16  { return utf16.Encode([]rune(s)) }
func hlString(s []uint16) string { return string(utf16.Decode(s)) }
func (d *hlDAT) transition(chars []uint16, state int) int {
	for _, c := range chars {
		p := state + int(c) + 1
		if p < 0 || p >= len(d.check) || int(d.check[p]) != state {
			return -1
		}
		state = int(d.base[p])
	}
	return state
}
func (d *hlDAT) output(state int) int {
	if state < 0 || state >= len(d.base) || int(d.check[state]) != state || d.base[state] >= 0 {
		return -1
	}
	return -int(d.base[state]) - 1
}
func (d *hlDAT) lookup(chars []uint16) int { return d.output(d.transition(chars, int(d.base[0]))) }
func (m *hlModel) nature(name string) int {
	for i, n := range m.natures {
		if n == name {
			return i
		}
	}
	panic("HanLP 未知词性 " + name)
}
func (m *hlModel) attr(name string, frequency int) *hlAttr {
	return &hlAttr{frequency, []hlLabel{{m.nature(name), frequency}}}
}
func (m *hlModel) vertex(chars []uint16, attr *hlAttr, id int, compile bool) *hlVertex {
	v := &hlVertex{real: chars, attr: attr, id: id}
	if !compile || len(attr.labels) != 1 {
		return v
	}
	nature := m.natures[attr.labels[0].tag]
	key := ""
	replace := false
	switch {
	case strings.HasPrefix(nature, "nr"):
		key = "未##人"
	case strings.HasPrefix(nature, "ns"):
		key = "未##地"
	case nature == "nx":
		key = "未##专"
		if m.ids[key] < 0 {
			key = "未##串"
		}
	case strings.HasPrefix(nature, "nt") || nature == "nit":
		key = "未##团"
	case strings.HasPrefix(nature, "m"):
		key = "未##数"
		replace = true
	case strings.HasPrefix(nature, "x"):
		key = "未##串"
		replace = true
	case nature == "t":
		key = "未##时"
		replace = true
	}
	if key != "" {
		v.id = m.ids[key]
		if replace {
			v.attr = m.core.attrs[v.id]
		}
	}
	return v
}
func (m *hlModel) regularVertex(chars []uint16) *hlVertex {
	id := m.core.lookup(chars)
	if id < 0 {
		return m.vertex(chars, m.attr("n", 1), -1, true)
	}
	return m.vertex(chars, m.core.attrs[id], id, true)
}
func (net hlNet) add(line int, v *hlVertex) bool {
	for _, old := range net[line] {
		if len(old.real) == len(v.real) {
			return false
		}
	}
	net[line] = append(net[line], v)
	return true
}
func (m *hlModel) wordNet(chars []uint16) hlNet {
	net := make(hlNet, len(chars)+2)
	net[0] = []*hlVertex{m.vertex(hlChars(" "), m.attr("begin", m.total/10), m.ids["始##始"], false)}
	net[len(net)-1] = []*hlVertex{m.vertex(hlChars(" "), m.attr("end", m.total/10), m.ids["末##末"], false)}
	for begin := range chars {
		state := int(m.core.base[0])
		for end := begin; end < len(chars); end++ {
			state = m.core.transition(chars[end:end+1], state)
			if state < 0 {
				break
			}
			id := m.core.output(state)
			if id >= 0 {
				net.add(begin+1, m.vertex(chars[begin:end+1], m.core.attrs[id], id, true))
			}
		}
	}
	for i := 1; i < len(net); {
		if len(net[i]) > 0 {
			i += len(net[i][len(net[i])-1].real)
			continue
		}
		j := i + 1
		for j < len(net)-1 && (len(net[j]) == 0 || m.charType[chars[j-1]] == 11) {
			j++
		}
		start := i - 1
		end := j - 1
		kind := m.charType[chars[start]]
		flush := func(to int) {
			nature := "n"
			id := -1
			switch kind {
			case 5:
				nature = "nx"
				id = m.ids["未##串"]
			case 6, 17:
				nature = "w"
			case 9, 10, 11:
				nature = "m"
				id = m.ids["未##数"]
			}
			net.add(start+1, m.vertex(chars[start:to], m.attr(nature, m.oov), id, false))
		}
		for k := start + 1; k < end; k++ {
			next := m.charType[chars[k]]
			if next != kind {
				if kind == 9 && strings.ContainsRune("，,．.", rune(chars[k])) && k+1 < end && m.charType[chars[k+1]] == 9 {
					continue
				}
				flush(k)
				start = k
			}
			kind = next
		}
		flush(end)
		i = j
	}
	return net
}
func (m *hlModel) bigram(from, to int) int {
	if from < 0 {
		return -from
	}
	if to < 0 {
		return -to
	}
	start, end := int(m.start[from]), int(m.start[from+1])
	i := sort.Search(end-start, func(i int) bool { return int(m.pair[(start+i)*2]) >= to })
	if start+i < end && int(m.pair[(start+i)*2]) == to {
		return int(m.pair[(start+i)*2+1])
	}
	return 0
}
func (m *hlModel) viterbi(net hlNet) []*hlVertex {
	update := func(from, to *hlVertex) {
		weight := from.weight - math.Log(0.9*(m.myu*float64(m.bigram(from.id, to.id))/float64(from.attr.total+1)+1.0-m.myu)+0.09999999999999998*float64(to.attr.total)/float64(m.total))
		if to.from == nil || to.weight > weight {
			to.from = from
			to.weight = weight
		}
	}
	for _, node := range net[1] {
		update(net[0][0], node)
	}
	for i := 1; i < len(net)-1; i++ {
		for _, node := range net[i] {
			if node.from == nil {
				continue
			}
			for _, to := range net[i+len(node.real)] {
				update(node, to)
			}
		}
	}
	path := []*hlVertex{}
	for node := net[len(net)-1][0]; node != nil; node = node.from {
		path = append(path, node)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
func (m *hlModel) combineCustom(path []*hlVertex) []*hlVertex {
	combine := func(i, end int, attr *hlAttr) {
		if i+1 == end {
			path[i].attr = attr
			return
		}
		chars := []uint16{}
		for j := i; j < end; j++ {
			if path[j] != nil {
				chars = append(chars, path[j].real...)
				path[j] = nil
			}
		}
		path[i] = m.vertex(chars, attr, -attr.total, false)
	}
	for i := 1; i < len(path)-1; i++ {
		state := m.custom.transition(path[i].real, 1)
		if state < 0 {
			continue
		}
		id := m.custom.output(state)
		end := i + 1
		for j := i + 1; j < len(path)-1; j++ {
			state = m.custom.transition(path[j].real, state)
			if state < 0 {
				break
			}
			if next := m.custom.output(state); next >= 0 {
				id = next
				end = j + 1
			}
		}
		if id >= 0 {
			combine(i, end, m.custom.attrs[id])
			i = end - 1
		}
	}
	// 原业务每次 CustomDictionary.add 的两个词，按 BinTrie 在主词典合并后追加。
	for _, word := range []string{"主焦", "1/3焦"} {
		chars := hlChars(word)
		if m.custom.lookup(chars) >= 0 {
			continue
		}
		for i := 1; i < len(path)-1; i++ {
			if path[i] == nil {
				continue
			}
			joined := []uint16{}
			for j := i; j < len(path)-1; j++ {
				if path[j] == nil {
					continue
				}
				joined = append(joined, path[j].real...)
				if len(joined) > len(chars) {
					break
				}
				if hlString(joined) == word {
					combine(i, j+1, m.attr("nz", 1))
					i = j
					break
				}
			}
		}
	}
	out := []*hlVertex{}
	for _, v := range path {
		if v != nil {
			out = append(out, v)
		}
	}
	return out
}

func (m *hlModel) role(v *hlVertex) []hlLabel {
	if id := m.nr.lookup(v.real); id >= 0 {
		return m.nr.roles[id]
	}
	nature := m.natures[v.attr.labels[0].tag]
	if nature == "nr" && v.attr.total <= 1000 && len(v.real) == 2 {
		return []hlLabel{{5, 1}, {11, 2}}
	}
	if nature == "nnt" {
		return []hlLabel{{5, 1}, {6, 1}}
	}
	return []hlLabel{{14, m.nrTotal[14]}}
}
func hlNetSize(net hlNet) int {
	size := 0
	for _, row := range net {
		size += len(row)
	}
	return size
}
func (opt hlNet) insert(line int, v *hlVertex, all hlNet) {
	if !opt.add(line, v) {
		return
	}
	start := line - 5
	if start < 0 {
		start = 0
	}
	for l := line - 1; l > start; l-- {
		if len(all[l]) > len(opt[l]) {
			for _, pre := range all[l] {
				if len(pre.real)+l == line {
					opt[l] = append(opt[l], pre)
				}
			}
		}
	}
	l := line + len(v.real)
	if len(opt[l]) == 0 && len(all[l]) > 0 {
		opt[l] = all[l]
	}
}
func (m *hlModel) recognizeNames(path []*hlVertex, opt, all hlNet) {
	tags := make([]int, len(path))
	tags[0] = 6
	pre := 6 // EnumItem(A,K) 的 TreeMap 首项为 K。
	for i := 1; i < len(path); i++ {
		cost := math.MaxFloat64
		tag := pre
		for _, label := range m.role(path[i]) {
			now := m.nrTransition[pre][label.tag] - math.Log((float64(label.frequency)+1e-8)/float64(m.nrTotal[label.tag]))
			if cost > now {
				cost = now
				tag = label.tag
			}
		}
		tags[i] = tag
		pre = tag
	}
	// U/V 复合角色拆字仅作用于模式匹配，原粗分路径保持不变。
	letters := "BCDEFGKLMUVXYZAS"
	pattern := []byte{}
	words := []*hlVertex{}
	previous := 14
	for i, tag := range tags {
		v := path[i]
		switch tag {
		case 9:
			pattern = append(pattern, 'K', 'B')
			words = append(words, m.regularVertex(v.real[:len(v.real)-1]), m.regularVertex(v.real[len(v.real)-1:]))
			previous = 0
		case 10:
			letter := byte('D')
			if previous == 0 {
				letter = 'E'
			}
			pattern = append(pattern, letter, 'L')
			words = append(words, m.regularVertex(v.real[:1]), m.regularVertex(v.real[1:]))
		default:
			pattern = append(pattern, letters[tag])
			words = append(words, v)
			previous = tag
		}
	}
	offsets := make([]int, len(words))
	for i := 1; i < len(words); i++ {
		offsets[i] = offsets[i-1] + len(words[i-1].real)
	}
	patterns := []string{"BBCD", "BBE", "BBZ", "BCD", "BEE", "BE", "BC", "BEC", "BG", "DG", "EG", "BZ", "EE", "FE", "FC", "FB", "FG", "Y", "XD"}
	for end := 1; end <= len(pattern); end++ {
		for _, p := range patterns {
			begin := end - len(p)
			if begin < 0 || string(pattern[begin:end]) != p {
				continue
			}
			name := []uint16{}
			for _, v := range words[begin:end] {
				name = append(name, v.real...)
			}
			if p == "BCD" && name[0] == name[2] {
				continue
			}
			bad := false
			if id := m.nr.lookup(name); id >= 0 {
				for _, label := range m.nr.roles[id] {
					if label.tag == 14 {
						bad = true
						break
					}
				}
			}
			if !bad {
				opt.insert(offsets[begin], m.vertex(name, m.attr("nr", 100), m.ids["未##人"], false), all)
			}
		}
	}
}
func (m *hlModel) recognizeTranslated(path []*hlVertex, opt, all hlNet) {
	name := []uint16{}
	appendTimes, line, activeLine := 0, 1, 1
	for _, v := range path[1:] {
		nature := m.natures[v.attr.labels[0].tag]
		if appendTimes > 0 {
			if nature != "nrf" && m.nrf.lookup(v.real) < 0 {
				if appendTimes > 1 {
					opt.insert(activeLine, m.vertex(name, m.attr("nrf", 1000), m.ids["未##人"], false), all)
				}
				name = nil
				appendTimes = 0
			} else {
				name = append(name, v.real...)
				appendTimes++
			}
		} else if nature == "nrf" {
			name = append(name, v.real...)
			appendTimes++
			activeLine = line
		}
		line += len(v.real)
	}
}

func hanlpNative(s string) []string {
	if s == "" {
		return []string{}
	}
	m := loadHanlp()
	chars := hlChars(s)
	if m.normalization {
		for i, c := range chars {
			chars[i] = m.charTable[c]
		}
	}
	all := m.wordNet(chars)
	path := m.combineCustom(m.viterbi(all))
	opt := make(hlNet, len(all))
	line := 0
	for _, v := range path {
		opt[line] = append(opt[line], v)
		line += len(v.real)
	}
	before := hlNetSize(opt)
	m.recognizeNames(path, opt, all)
	m.recognizeTranslated(path, opt, all)
	if hlNetSize(opt) != before {
		path = m.viterbi(opt)
	}
	out := make([]string, 0, len(path)-2)
	for _, v := range path[1 : len(path)-1] {
		out = append(out, hlString(v.real))
	}
	return out
}
