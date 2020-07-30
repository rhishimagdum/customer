package main

func main() {
	a := App{}
	a.Initialize(CustomerImplCassandra{}, KafkaMessanger{})
	a.Run(":80")
}
