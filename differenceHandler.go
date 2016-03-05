package consulKvMigrator

type differenceHandler interface {
	handle(differences differences) error
	setNext(differenceHandler differenceHandler)
}

type chainableDifferenceHandler struct {
	nextHandler differenceHandler
}

func (handler *chainableDifferenceHandler) next(differences differences) error {
	if handler.nextHandler != nil {
		return handler.nextHandler.handle(differences)
	}

	return nil
}

func (handler *chainableDifferenceHandler) setNext(nextHandler differenceHandler) {
	handler.nextHandler = nextHandler
}

func buildHandlers(consulKvClient consulKvClient, prompt bool, dryRun bool) differenceHandler {
	handlers := []differenceHandler{&reportHandler{dryRun: dryRun}}
	if prompt {
		handlers = append(handlers, &promptingHandler{dryRun: dryRun})
	}
	if !dryRun {
		handlers = append(handlers, &consulHandler{consulKvClient: consulKvClient})
	}

	return chainHandlers(handlers)
}

func chainHandlers(handlers []differenceHandler) differenceHandler {
	head := handlers[0]
	current := head
	idx := 1
	for idx < len(handlers) {
		next := handlers[idx]
		current.setNext(next)
		idx++
		current = next
	}
	return head
}
