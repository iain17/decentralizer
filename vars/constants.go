package vars

import (
	"github.com/c2h5oh/datasize"
	"github.com/iain17/framed"
	"time"
)

const MAX_DISCOVERED_PEERS = 30
const MAX_BOOTSTRAP_SHARE = 80
var MIN_CONNECTED_PEERS = 2//40
const DELIMITER_ADDR = ";;"
var EXPIRE_TIME_SESSION = 30 * time.Minute
var DIFF_DIFFERENCE_ACCEPTANCE = 1 * time.Minute
const MAX_SESSIONS = 1000
const MAX_CONTACTS = 1000
const EXPIRE_TIME_CONTACT = 10800//3 hours
const SEND_DIRECT_MESSAGE = "/decentralizer/dm/1.0.0/sent"
const ADDRESS_BOOK_FILE = "addressbook.dat"
const SESSIONS_FILE = "sessions.dat"
const PUBLISHER_DEFINITION_FILE = "publisherDefinition.dat"
const MAX_SIZE = int64(10 * datasize.MB)
const MAX_IGNORE = 4096//If a peer isn't using our protocol. max ignore
const MAX_UNMARSHAL_CACHE = 500//If a peer isn't using our protocol. max ignore
const CONCURRENT_SESSION_REQUEST = 10
const GET_SESSION_REQ = "/decentralizer/sessions/1.0.0/get"
const MAX_SESSION_SEARCHES = 10
const MESSAGE_DEADLINE = time.Minute * 10
const DHT_PEER_KEY_TYPE = "dPeer"
const DHT_SESSIONS_KEY_TYPE = "sessions"
const DHT_PUBLISHER_KEY_TYPE = "publisher"
const BOOTSTRAP_FILE = "bootstrap.dat"
const NUM_MATCHMAKING_SLOTS = int(100)
const MAX_IDLE_TIME = 5 * time.Minute//Ignored in daemon mode

var FILE_EXPIRE = time.Hour * 1

func init() {
	framed.MAX_SIZE = MAX_SIZE
}