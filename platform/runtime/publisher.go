package runtime

// Publisher defines metric publisher contract
type Publisher interface {
	Publish(rs *Stats)
}
