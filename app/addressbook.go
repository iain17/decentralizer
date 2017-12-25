package app

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	"gx/ipfs/QmT6n4mspWYEya864BhCUJEgyxiRfmiSY9ruQwTUNpRKaM/protobuf/proto"
	inet "gx/ipfs/QmU4vCDZTPLDqSDKguWbHCiUe46mZUtmM2g2suBZ9NE8ko/go-libp2p-net"
	Peer "gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"github.com/Pallinder/go-randomdata"
	pstore "gx/ipfs/QmYijbtjCxFEjSXaudaQAUz3LN5VKLssm8WCUsRoqzXmQR/go-libp2p-peerstore"
	ma "gx/ipfs/QmW8s4zTsUoX1Q6CeYxVKPyqSKbF7H1YDUyTostBtZ8DaG/go-multiaddr"
	"time"
	"github.com/iain17/framed"
)

func (d *Decentralizer) initAddressbook() {
	d.i.PeerHost.SetStreamHandler(GET_PEER_REQ, d.getPeerResponse)
	d.downloadPeers()
	d.saveSelf()
	go d.connectPreviousPeers()
	d.provideSelf()
	d.cron.AddFunc("30 * * * * *", d.uploadPeers)
	d.cron.AddFunc("* 5 * * * *", d.provideSelf)
}

func (d *Decentralizer) downloadPeers() {
	data, err := configPath.QueryCacheFolder().ReadFile(ADDRESS_BOOK_FILE)
	if err != nil {
		//logger.Warningf("Could not restore address book: %v", err)
		return
	}
	var addressbook pb.DNAddressbook
	err = proto.Unmarshal(data, &addressbook)
	if err != nil {
		logger.Warningf("Could not restore address book: %v", err)
		return
	}
	for _, peer := range addressbook.Peers {
		err = d.peers.Upsert(peer)
		if err != nil {
			logger.Warningf("Error saving peer: %s", peer.PId)
			continue
		}
	}
	logger.Info("Restored address book")
}

func (d *Decentralizer) provideSelf() {
	peer, err := d.FindByPeerId("self")
	if err != nil {
		logger.Warningf("Could not provide self: %v", err)
		return
	}
	d.b.Provide(getDecentralizedIdKey(peer.DnId))
	logger.Debug("Provided self")
}

func (d *Decentralizer) uploadPeers() {
	if !d.addressBookChanged {
		return
	}
	peers, err := d.peers.FindAll()
	if err != nil {
		logger.Warningf("Could not save address book: %v", err)
		return
	}
	data, err := proto.Marshal(&pb.DNAddressbook{
		Peers: peers,
	})
	if err != nil {
		logger.Warningf("Could not save address book: %v", err)
		return
	}
	err = configPath.QueryCacheFolder().WriteFile(ADDRESS_BOOK_FILE, data)
	if err != nil {
		logger.Warningf("Could not save address book: %v", err)
		return
	}
	d.addressBookChanged = true
	logger.Info("Saved address book")
}

//Connect to our previous peers
func (d *Decentralizer) connectPreviousPeers() error {
	for {
		lenPeers := len(d.i.PeerHost.Network().Peers())
		if lenPeers >= 3 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	logger.Info("Connecting to previous peers...")
	peers, err := d.peers.FindAll()
	if err != nil {
		return err
	}
	for _, peer := range peers {
		pId, err := d.decodePeerId(peer.PId)
		if err != nil {
			continue
		}
		var addrs []ma.Multiaddr
		for _, rawAddr := range peer.Addrs {
			addr, err := ma.NewMultiaddr(rawAddr)
			if err != nil {
				continue
			}
			addrs = append(addrs, addr)
		}
		err = d.i.PeerHost.Connect(d.i.Context(), pstore.PeerInfo{
			ID: pId,
			Addrs: addrs,
		})
		if err != nil {
			logger.Warning(err)
		}
	}
	return nil
}

//Save ourself at least in the address book.
func (d *Decentralizer) saveSelf() error {
	self, err := d.peers.FindByPeerId(d.i.Identity.Pretty())
	if err != nil || self == nil {
		//Add self
		err = d.UpsertPeer(d.i.Identity.Pretty(), map[string]string{
			"name": randomdata.SillyName(),
		})
		if err != nil {
			return err
		}
		d.uploadPeers()
	}
	return nil
}

func (d *Decentralizer) UpsertPeer(pId string, details map[string]string) error {
	err := d.peers.Upsert(&pb.Peer{
		PId:     pId,
		Details: details,
	})
	d.addressBookChanged = true
	return err
}

func (d *Decentralizer) GetPeersByDetails(key, value string) ([]*pb.Peer, error) {
	return d.peers.FindByDetails(key, value)
}

func (d *Decentralizer) GetPeers() ([]*pb.Peer, error) {
	return d.peers.FindAll()
}

func (d *Decentralizer) FindByPeerId(peerId string) (p *pb.Peer, err error) {
	p, err = d.peers.FindByPeerId(peerId)
	if err != nil {
		var id Peer.ID
		id, err = d.decodePeerId(peerId)
		if err != nil {
			return nil, err
		}
		p, err = d.getPeerRequest(id)
		if err != nil {
			return nil, err
		}
		d.peers.Upsert(p)
	}
	return p, err
}

func (d *Decentralizer) getPeerRequest(peer Peer.ID) (*pb.Peer, error) {
	stream, err := d.i.PeerHost.NewStream(d.i.Context(), peer, GET_PEER_REQ)
	if err != nil {
		return nil, err
	}
	stream.SetDeadline(time.Now().Add(300 * time.Millisecond))
	defer stream.Close()

	//Request
	reqData, err := proto.Marshal(&pb.DNPeerRequest{})
	if err != nil {
		return nil, err
	}
	err = framed.Write(stream, reqData)
	if err != nil {
		return nil, err
	}

	//Response
	resData, err := framed.Read(stream)
	if err != nil {
		return nil, err
	}
	var response pb.DNPeerResponse
	err = proto.Unmarshal(resData, &response)
	if err != nil {
		return nil, err
	}

	//Save addr so we can quickly connect to our contacts
	if response.Peer != nil {
		info := d.i.Peerstore.PeerInfo(peer)
		response.Peer.Addrs = []string{}
		for _, addr := range info.Addrs {
			response.Peer.Addrs = append(response.Peer.Addrs, addr.String())
		}
	}
	return response.Peer, nil
}

func (d *Decentralizer) FindByDecentralizedId(decentralizedId uint64) (*pb.Peer, error) {
	peer, err := d.peers.FindByDecentralizedId(decentralizedId)
	if err != nil || peer == nil {
		peerId, err := d.resolveDecentralizedId(decentralizedId)
		if err != nil {
			return nil, err
		}
		return d.FindByPeerId(peerId.Pretty())
	}
	return peer, err
}

func (d *Decentralizer) getPeerResponse(stream inet.Stream) {
	reqData, err := framed.Read(stream)
	if err != nil {
		logger.Error(err)
		return
	}
	var request pb.DNPeerRequest
	err = proto.Unmarshal(reqData, &request)
	if err != nil {
		logger.Error(err)
		return
	}
	peer, err := d.peers.FindByPeerId(d.i.Identity.Pretty())
	if err != nil {
		logger.Error(err)
		return
	}

	//Response
	response, err := proto.Marshal(&pb.DNPeerResponse{
		Peer: peer,
	})
	if err != nil {
		logger.Error(err)
		return
	}
	err = framed.Write(stream, response)
	if err != nil {
		logger.Error(err)
		return
	}
}
