package logger

type Output interface {
	Print(int, string) error
}

var outputs []Output

func AddOutput(output Output) {
	outputs = append(outputs, output)
}
