package app

import (
	"testing"
	"context"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/getlantern/testify/assert"
	"github.com/iain17/decentralizer/pb"
)

//In this test case: app1 wants to retrieve publisher file from the publisher definition
func TestDecentralizer_GetPublisherFile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx,2)
	app1 := fakeNew(nodes[0], false)
	assert.NotNil(t, app1)
	expected := []byte("Ok")

	//Mocked publisher update
	app1.publisherUpdate = &pb.PublisherUpdate{}
	app1.publisherDefinition = &pb.PublisherDefinition{
		Status: true,
		Files: map[string][]byte {
			"test.txt": expected,
		},
	}

	data, err := app1.GetPublisherFile("test.txt")
	assert.NoError(t, err)
	assert.Equal(t, expected, data)
}

//In this test case: app1 wants to retrieve publisher file that is on IPFS. Which app2 has published.
func TestDecentralizer_GetPublisherFile2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx,2)
	app1 := fakeNew(nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(nodes[1], false)
	assert.NotNil(t, app2)

	const filename = "test.txt"
	expected := []byte("Ok")
	ipfsPath, err := app2.SavePeerFile(filename, expected)
	assert.NoError(t, err)

	publisherUpdate := &pb.PublisherUpdate{}
	definition := &pb.PublisherDefinition{
		Status: true,
		Links: map[string]string {
			filename: ipfsPath,
		},
	}
	//
	////Mocked publisher update
	app1.publisherUpdate = publisherUpdate
	app1.publisherDefinition = definition
	app2.publisherUpdate = publisherUpdate
	app2.publisherDefinition = definition

	data, err := app1.GetPublisherFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, expected, data)
}
