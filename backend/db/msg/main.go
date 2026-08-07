package msg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type EventData string

type Part interface {
	render(*strings.Builder)
}

type text string

func (t text) render(sb *strings.Builder) {
	sb.WriteString(string(t))
}

func Text(s string) Part {
	return text(s)
}

func Int(n int) Part {
	return text(strconv.Itoa(n))
}

type wrapped struct {
	tag      string
	children []Part
}

func (w wrapped) render(sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("[%s]", w.tag))
	for _, c := range w.children {
		c.render(sb)
	}
	sb.WriteString(fmt.Sprintf("[/%s]", w.tag))
}

type group struct {
	children []Part
}

func (g group) render(sb *strings.Builder) {
	for _, c := range g.children {
		c.render(sb)
	}
}

func Bold(n ...Part) Part {
	return wrapped{"b", n}
}

func Player(id uuid.UUID) Part {
	return wrapped{"player", []Part{text(id.String())}}
}
func Icon(id string) Part {
	return wrapped{"icon", []Part{text(id)}}
}

func Concat(n ...Part) Part {
	return group{n}
}

type MessageBuilder struct {
	message []Part
}

type MessageTag struct {
	tag string
}

func Msg() *MessageBuilder {
	return &MessageBuilder{}
}

func (m *MessageBuilder) Add(s Part) *MessageBuilder {
	m.message = append(m.message, s)
	return m
}

func (m *MessageBuilder) Player(id uuid.UUID) *MessageBuilder {
	m.Add(Player(id))
	return m
}

func (m *MessageBuilder) Icon(v string) *MessageBuilder {
	m.Add(Icon(v))
	return m
}

func (m *MessageBuilder) Text(s string) *MessageBuilder {
	m.Add(Text(s))
	return m
}

func (m *MessageBuilder) Bold(s Part) *MessageBuilder {
	m.Add(Bold(s))
	return m
}

func (m *MessageBuilder) String() EventData {
	var sb strings.Builder
	for _, n := range m.message {
		n.render(&sb)
	}
	return EventData(sb.String())
}
