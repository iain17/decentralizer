package app

import "github.com/c2h5oh/datasize"

const MAX_DISCOVERED_PEERS = 10
const MIN_DISCOVERED_PEERS = 1//min 3. ipfs requires min of 3.
const DELIMITER_ADDR = ";;"
const EXPIRE_TIME_SESSION = 30
const MAX_SESSIONS = 1000
const GET_SESSION_REQ = "/decentralizer/1.0.0/get"
const MAX_SIZE = int64(3 * datasize.MB)