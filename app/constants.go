package app

import "github.com/c2h5oh/datasize"

const MAX_DISCOVERED_PEERS = 10
const MIN_CONNECTED_PEERS = 1//min 3. ipfs requires min of 3.
const DELIMITER_ADDR = ";;"
const EXPIRE_TIME_SESSION = 30
const MAX_SESSIONS = 1000
const MAX_CONTACTS = 1000
const EXPIRE_TIME_CONTACT = 3600
const GET_PEER_REQ = "/decentralizer/peers/1.0.0/get"
const GET_SESSION_REQ = "/decentralizer/sessions/1.0.0/get"
const SEND_DIRECT_MESSAGE = "/decentralizer/dm/1.0.0/sent"
const MAX_SIZE = int64(3 * datasize.MB)