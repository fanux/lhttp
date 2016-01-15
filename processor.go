package lhttp

var processorMap map[string]HandlerCallbacks

func Regist(command string, p HandlerCallbacks) {
	processorMap[command] = p
	//log.Print("regist processor", p)
}

func getProcessor(command string) HandlerCallbacks {
	//log.Print("get processor:", processorMap[command], " command:", command)
	p, ok := processorMap[command]
	if ok {
		return p
	} else {
		return nil
	}
}

func init() {
	processorMap = make(map[string]HandlerCallbacks, ProcessorMax)
}
