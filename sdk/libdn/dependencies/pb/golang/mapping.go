package np

import (
	"errors"
	"gonp/np/handlers"
	"gonp/np/structs"
	"net"
)

var NoHandlerFound = errors.New("No handler found")

func HandleMessage(conn net.Conn, connection_data *structs.ConnData, packet_data *structs.PacketData, np_server *NPServer) error {
	switch packet_data.Header.Type {
		 case 1000:
			return handlers.RPCHelloRequest(conn, connection_data, packet_data, np_server.Query_server.Listener)
		 case 1001:
			return handlers.RPCHelloReply(conn, connection_data, packet_data, np_server.Query_server.Listener)
		}

	return NoHandlerFound
}