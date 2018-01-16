package discovery

import (
	"github.com/shibukawa/configdir"
	"time"
	"github.com/iain17/framed"
	"github.com/c2h5oh/datasize"
)

const CONCCURENT_NEW_CONNECTION = 50
const BACKLOG_NEW_CONNECTION = 100
const HEARTBEAT_DELAY = 30
const PORT_RANGE = 10
const IRC_SERVER = "chat.freenode.net:6697"
const IRC_PREFIX = "JOIN US:"
const NET_TABLE_FILE = "peers.data"
const SERVICE = "_decentralizer._tcp"
const DEADLINE_DURATION = 1 * time.Second
var configPath = configdir.New("ECorp", "Decentralizer")

func init()  {
	framed.MAX_SIZE = int64(5 * datasize.MB)
}