package main

type Messanger interface {
	PublishMessage(m string)
	ConsumeMessage()
}
