package irc

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

// CommandFilter filters events based on the IRC command of the event
type CommandFilter struct {
	Command string
}

// Match this filter, against incoming events.
func (cf CommandFilter) Match(ev *Event) bool {
	if cf.Command == "*" || cf.Command == "" {
		// Match all events
		return true
	}
	return cf.Command == ev.Command
}