package setting

import "os"

func Neo4jConnectionString() string {
	return os.Getenv("NEO4J_CONNECTION_STRING")
}
func Neo4jUser() string {
	return os.Getenv("NEO4J_USER")
}
func Neo4jPassword() string {
	return os.Getenv("NEO4J_PASSWORD")
}
