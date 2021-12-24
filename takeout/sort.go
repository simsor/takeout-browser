package takeout

// ByTimestamp implements the Sort interface, sorting from lowest timestamp to highest
type ByTimestamp []Media

func (m ByTimestamp) Len() int {
	return len(m)
}

func (m ByTimestamp) Less(i, j int) bool {
	return m[i].Taken < m[j].Taken
}

func (m ByTimestamp) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// ByReverseTimestamp implements the Sort interface, sorting from highest timestamp to lowest
type ByReverseTimestamp []Media

func (m ByReverseTimestamp) Len() int {
	return len(m)
}

func (m ByReverseTimestamp) Less(i, j int) bool {
	return m[i].Taken > m[j].Taken
}

func (m ByReverseTimestamp) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
