package irc

import "regexp"

// A Filter is ways to match Events from the server to some matching pattern,
// so that specific events can trigger specific actions (or not trigger actions)
type Filter interface {
	Match(ev *Event) bool
}

// HandlerFunc preforms some action, based on the Event given, and respsonds
// using the IRC client c
type HandlerFunc func(ev *Event, c IRC)

// Handler filters events to be acted on by a HandlerFunc
type Handler struct {
	Filters []Filter
	Handler HandlerFunc
}

// CombinedFilter is a filter for matching events against multiple Filters
type CombinedFilter struct {
	fs []Filter
}

// NewCombinedFilter returns a new CombinedFilter
func NewCombinedFilter(fs ...Filter) *CombinedFilter {
	return &CombinedFilter{
		fs: fs,
	}
}

// Match matches the filter against incoming events
func (cf *CombinedFilter) Match(ev *Event) bool {
	for _, f := range cf.fs {
		if val := f.Match(ev); !val {
			return false
		}
	}

	return true
}

// CommandFilter filters events based on the IRC command of the event
type CommandFilter string

// Match events for CommandFilter
func (cf CommandFilter) Match(ev *Event) bool {
	if cf == "*" || cf == "" {
		// Match all events
		return true
	}
	return cf == CommandFilter(ev.Command)
}

// RegExpFilterType is the Event paremter we should match on
type RegExpFilterType int

const (
	// RegExpFilterCommand filters against Event.Command
	RegExpFilterCommand = iota

	// RegExpFilterPrefix filters against Event.Prefix
	RegExpFilterPrefix = iota

	// RegExpFilterParameters filters against (all) event.Parameters
	RegExpFilterParameters = iota
)

// RegExpFilter filters message content by Regular Expression
type RegExpFilter struct {
	Param      RegExpFilterType
	Expression regexp.Regexp
}

// Match Events for RegExpFilter
func (ref *RegExpFilter) Match(ev *Event) bool {
	find := [][]byte{}

	switch ref.Param {
	case RegExpFilterCommand:
		find = append(find, []byte(ev.Command))
	case RegExpFilterPrefix:
		find = append(find, []byte(ev.Prefix))
	case RegExpFilterParameters:
		for _, p := range ev.Parameters {
			find = append(find, []byte(p))
		}
	}

	for _, f := range find {
		if ref.Expression.Find(f) != nil {
			return true
		}
	}

	return false
}

// ChannelFilter filters events by the channel they were sent
type ChannelFilter string

// Match Events for ChannelFilter
func (cf ChannelFilter) Match(ev *Event) bool {
	if len(ev.Parameters) == 0 {
		return false
	}

	return ChannelFilter(ev.Parameters[0]) == cf
}
