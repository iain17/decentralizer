package app

import (
	"github.com/c2h5oh/datasize"
	"github.com/iain17/framed"
	"time"
)

const USE_OWN_BOOTSTRAPPING = true//If set to false. We join the public IPFS network.
const MAX_DISCOVERED_PEERS = 10
var MIN_CONNECTED_PEERS = 1//40
const DELIMITER_ADDR = ";;"
var EXPIRE_TIME_SESSION = 3 * time.Hour
const MAX_SESSIONS = 1000
const MAX_CONTACTS = 1000
const EXPIRE_TIME_CONTACT = 10800//3 hours
const SEND_DIRECT_MESSAGE = "/decentralizer/dm/1.0.0/sent"
const ADDRESS_BOOK_FILE = "addressbook.dat"
const PUBLISHER_DEFINITION_FILE = "publisherDefinition.dat"
const MAX_SIZE = int64(10 * datasize.MB)
const MAX_IGNORE = 4096//If a peer isn't using our protocol. max ignore
const CONCURRENT_SESSION_REQUEST = 100
const MAX_SESSION_SEARCHES = 10
const MESSAGE_DEADLINE = time.Minute * 10
const DHT_PEER_KEY_TYPE = "dPeer"
const DHT_SESSIONS_KEY_TYPE = "sessions"
var FILE_EXPIRE = time.Hour * 1

func init() {
	framed.MAX_SIZE = MAX_SIZE
}