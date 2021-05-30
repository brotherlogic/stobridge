package main

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	return s
}
