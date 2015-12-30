package lhttp

var (
	ProcessorMax = 40
)

var processorMap map[string]HandlerCallbacks

func Regist(command string, p HandlerCallbacks) {
	processorMap[command] = p
	//log.Print("regist processor", p)
}

func getProcessor(command string) HandlerCallbacks {
	//log.Print("get processor:", processorMap[command], " command:", command)
	return processorMap[command]
}

func init() {
	processorMap = make(map[string]HandlerCallbacks, ProcessorMax)
}
