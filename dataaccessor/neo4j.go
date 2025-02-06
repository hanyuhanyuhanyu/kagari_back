package dataaccessor

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ConnectionInfo struct {
	ConnectionString string
	User             string
	Password         string
}

type connectionHandler struct {
	driver neo4j.DriverWithContext
}

func (ch *connectionHandler) Open(ctx context.Context, conInfo ConnectionInfo) (neo4j.DriverWithContext, error) {
	if ch.driver != nil {
		return ch.driver, nil
	}
	driver, err := neo4j.NewDriverWithContext(conInfo.ConnectionString, neo4j.BasicAuth(conInfo.User, conInfo.Password, ""))
	if err != nil {
		return nil, err
	}
	ch.driver = driver
	return ch.driver, nil
}
func (ch *connectionHandler) Close(ctx context.Context) error {
	if ch.driver == nil {
		return nil
	}
	return ch.driver.Close(ctx)
}
func WithNeo4jConnection(ctx context.Context, conInfo ConnectionInfo, f func(neo4cDriver neo4j.DriverWithContext)) error {
	ch := connectionHandler{}
	_, err := ch.Open(ctx, conInfo)
	if err != nil {
		return err
	}
	defer ch.Close(ctx)
	f(ch.driver)
	return nil
}
