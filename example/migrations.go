package main

import (
	"github.com/hinshun/pge"
	"github.com/hinshun/pge/example/sql"
)

var (
	SchemaMigrations = []pge.Migration{
		{
			Name: "Create initial schema",
			Queries: []pge.Query{
				sql.CreateTableCustomers,
			},
		},
		{
			Name: "Add product and transactions",
			Queries: []pge.Query{
				sql.CreateTableProducts,
				sql.CreateTableTransactions,
				sql.CreateTableProductTransactions,
			},
		},
	}
)
