package producer

type Producer interface {
	Produce(keySize uint)
}
