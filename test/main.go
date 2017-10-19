package main

import (
	"github.com/hashicorp/go-memdb"
	"fmt"
)

func main() {
	// Create a sample struct
	type Person struct {
		Email string
		Name  string
		Age   int
		Details map[string]string
	}

	/*
		Owner   *DPeer            `protobuf:"bytes,1,opt,name=owner" json:"owner,omitempty"`
		Type    uint32            `protobuf:"varint,2,opt,name=type" json:"type,omitempty"`
		Name    string            `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
		Address uint64            `protobuf:"varint,4,opt,name=address" json:"address,omitempty"`
		Port    uint32            `protobuf:"varint,5,opt,name=port" json:"port,omitempty"`
		Details map[string]string `protobuf:"bytes,6,rep,name=details" json:"details,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	 */


	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"person": &memdb.TableSchema{
				Name: "person",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"details": &memdb.IndexSchema{
						Name:    "details",
						Unique:  false,
						Indexer: &memdb.StringMapFieldIndex{Field: "Details"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	// Create a write transaction
	txn := db.Txn(true)

	// Insert a new person
	p := &Person{"joe@aol.com", "Joe", 30, map[string]string{
		"cool": "nope",
		"id": "000",
	}}
	if err := txn.Insert("person", p); err != nil {
		panic(err)
	}

	p = &Person{"iain@aol.com", "Iain", 23, map[string]string{
		"cool": "yes",
		"id": "0101",
	}}
	if err := txn.Insert("person", p); err != nil {
		panic(err)
	}

	p = &Person{"iain@aol.com", "Iain", 23, map[string]string{
		"cool": "yes",
		"id": "123",
	}}
	if err := txn.Insert("person", p); err != nil {
		panic(err)
	}

	// Commit the transaction
	txn.Commit()

	// Create read-only transaction
	txn = db.Txn(false)
	defer txn.Abort()

	// Lookup by email
	raw, err := txn.First("person", "details", "cool", "yes")
	if err != nil {
		panic(err)
	}

	// Say hi!
	fmt.Printf("Hello %s!", raw.(*Person).Name)

}