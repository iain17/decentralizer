package app

import "github.com/c2h5oh/datasize"

const USE_OWN_BOOTSTRAPPING = false//If set to false. We join the public IPFS network.
const MAX_DISCOVERED_PEERS = 40
const MIN_CONNECTED_PEERS = 40//40
const DELIMITER_ADDR = ";;"
const EXPIRE_TIME_SESSION = 120
const MAX_SESSIONS = 1000
const MAX_CONTACTS = 1000
const EXPIRE_TIME_CONTACT = 10800//3 hours
const GET_PEER_REQ = "/decentralizer/peers/1.0.0/get"
const GET_SESSION_REQ = "/decentralizer/sessions/1.0.0/get"
const SEND_DIRECT_MESSAGE = "/decentralizer/dm/1.0.0/sent"
const ADDRESS_BOOK_FILE = "addressbook.dat"
const PUBLISHER_TOPIC_FILES = "/decentralizer/publisher/files"
const MAX_SIZE = int64(3 * datasize.MB)