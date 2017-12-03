<?php
$script = new DataScript();
$script->read('messages.json');

generate_file($script->data, 'cpp-h.tpl', 'c++/MessageDefinition.h');
generate_file($script->data, 'cpp-c.tpl', 'c++/MessageDefinition.cpp');
generate_file($script->data, 'go.tpl', 'golang/reply/reply.go');
generate_file($script->data, 'gomap.tpl', 'golang/mapping.go');

function generate_file($messages, $template, $output)
{
	ob_start();
	include($template);
	$c = ob_get_contents();
	ob_end_clean();
	
	file_put_contents($output, $c);
}

class DataScript
{
	public $data = array();
	
	public function read($file)
	{
		$this->data = array();
		$messages = json_decode(file_get_contents($file), true);
		$type = 1000;
		foreach($messages as $message) {
			if(empty($message['message'])) {
				continue;
			}
			$message['name'] = "RPC".$message['message'];
			$message['type'] = $type;
			$this->data[] = $message;
			$type++;
		}
	}
	
	public function write()
	{
		$output = '';
	
		foreach ($this->data as $entity)
		{
			$output .= sprintf("{\n");
			
			foreach ($entity as $key => $value)
			{
				$output .= sprintf("\t\"%s\" \"%s\"\n", $key, $value);
			}
			
			$output .= sprintf("}\n");
		}
		
		return $output;
	}
}