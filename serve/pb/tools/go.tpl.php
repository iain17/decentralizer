// This file has been automatically generated.
package pb

import "reflect"

type MessageType int32
var types map[reflect.Type]MessageType

const (
<?php foreach ($messages as $message) {
	echo "	{$message['name']}	MessageType = {$message['type']}\n";
}
?>
)

func init() {
<?php foreach ($messages as $message) {
	echo "	registerType((*RPCMessage_{$message['message']})(nil), {$message['name']})\n";
}
?>
}

func registerType(x interface{}, msgType MessageType) {
	types[reflect.TypeOf(x)] = msgType
}