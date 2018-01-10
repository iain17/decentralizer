package discovery

import "github.com/shibukawa/configdir"

const CONCCURENT_NEW_CONNECTION = 10
const HEARTBEAT_DELAY = 30
const IRC_SERVER = "chat.freenode.net:6697"
const IRC_PREFIX = "JOIN US:"
const NET_TABLE_FILE = "peers.data"
const SERVICE = "_decentralizer._tcp"
var configPath = configdir.New("ECorp", "Decentralizer")