package models

type TaskType int

const (
	Permanent TaskType = iota + 1
	OneTime
)

func (taskType TaskType) String() string {
	// declare an array of strings
	// ... operator counts how many
	// items in the array (7)
	types := [...]string{
		"Permanent",
		"OneTime",
	}
	// → `day`: It's one of the
	// values of Weekday constants.
	// If the constant is Sunday,
	// then day is 0.
	// prevent panicking in case of
	// `day` is out of range of Weekday
	if taskType < Permanent || taskType > OneTime {
		return "Unknown"
	}
	// return the name of a Weekday
	// constant from the names array
	// above.
	return types[taskType]
}
