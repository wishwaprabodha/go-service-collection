package db

import (
	"fmt"
	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

func DetaConnection() (db *base.Base, err error) {
	conn, err := deta.New(deta.WithProjectKey("d0ifn5do_SRRFMLxr9zgsVwKDDSrHxB71LKE5RFT1"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return nil, err
	}
	database, err := base.New(conn, "note-server")
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return nil, err
	}
	return database, nil
}
