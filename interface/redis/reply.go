package redis

// Reply is the interface of redis serialization protocol message
type Reply interface {
	ToBytes() []byte
}

// SimpleStringReply is the reply of redis `+` command
type SimpleStringReply struct {
	Value string
}

// ToBytes implements Reply interface
func (reply *SimpleStringReply) ToBytes() []byte {
	return []byte(reply.Value)
}
