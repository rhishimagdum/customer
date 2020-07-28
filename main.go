package main

func main() {
	a := App{}
	a.Initialize("customer")
	a.Run(":80")
}
