package main


func main() { 
	server := NewAPI(":3000")
	server.Start() 
}