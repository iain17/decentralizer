package app

import (
	"context"
	"github.com/getlantern/testify/assert"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/iain17/decentralizer/pb"
	"testing"
	//"io/ioutil"
	//"os"
)

//In this test case: app1 wants to retrieve publisher file from the publisher definition
func TestDecentralizer_GetPublisherFile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx,2)
	app1 := fakeNew(ctx, nodes[0], false)
	assert.NotNil(t, app1)
	expected := []byte("Ok")

	//Mocked publisher update
	app1.publisherRecord = &pb.DNPublisherRecord{}
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
	app1 := fakeNew(ctx, nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(ctx, nodes[1], false)
	assert.NotNil(t, app2)

	const filename = "test.txt"
	expected := []byte("Ok")
	ipfsPath, err := app2.SavePeerFile(filename, expected)
	assert.NoError(t, err)
	app2.republishPeerFiles()//we could also wait 1 second.

	publisherUpdate := &pb.DNPublisherRecord{}
	definition := &pb.PublisherDefinition{
		Status: true,
		Links: map[string]string {
			filename: ipfsPath,
		},
	}

	////Mocked publisher update
	app1.publisherRecord = publisherUpdate
	app1.publisherDefinition = definition
	app2.publisherRecord = publisherUpdate
	app2.publisherDefinition = definition

	data, err := app1.GetPublisherFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, expected, data)
}

//func TestDecentralizer_AddPublisherFiles(t *testing.T) {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	nodes := ipfs.FakeNewIPFSNodes(ctx,10)
//	app1 := fakeNew(ctx, nodes[0], false)
//	assert.NotNil(t, app1)
//
//	const basedir = "./ipfs_test/"
//
//	os.RemoveAll(basedir)
//	os.MkdirAll(basedir, 0777)
//	ioutil.WriteFile(basedir+"hello.txt", []byte("this works"), 0777)
//	ioutil.WriteFile(basedir+"ok.txt", []byte("Nope!"), 0777)
//
//	err := app1.AddPublisherFiles(basedir)
//	os.RemoveAll(basedir)
//	assert.NoError(t, err)
//}