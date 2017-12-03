<?php
$returnMessages = array(1005, 1006,1007,1010, 1011, 1012,1022, 1111, 1112, 1113, 1211, 1215, 1217, 1302, 1304, 1306, 1308, 1310, 1311, 2001);

?>
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
	<?php
	foreach ($messages as $message)
	{?>
	<?= in_array($message['type'], $returnMessages) ? '//' : '' ?> case <?= $message['type'] ?>:
	<?= in_array($message['type'], $returnMessages) ? '//' : '' ?>		return handlers.<?= $message['name'] ?>(conn, connection_data, packet_data, np_server.Query_server.Listener)
	<?php
	}
	?>
	}

	return NoHandlerFound
}