package cli

type handlers interface {
	showHandlers(outputter Outputter)
}

type handlersForGroup struct {
	handlers []runnableHandler
}

var _ handlers = new(handlersForGroup)

func (hfg *handlersForGroup) getHandlers() []runnableHandler {
	return hfg.handlers
}

func (hfg *handlersForGroup) showHandlers(outputter Outputter) {
	parts := hfg.splitCommandsDialogsAndGroups()
	var lgs = lineGroups(parts)
	outputter(lgs.String())
}

func (hfg *handlersForGroup) splitCommandsDialogsAndGroups() []*lineGroup {
	groups := make([]runnableHandler, 0)
	commands := make([]runnableHandler, 0)
	dialogs := make([]runnableHandler, 0)
	for _, handler := range hfg.handlers {
		switch specificHandler := handler.(type) {
		case *groupContextRunnable:
			groups = append(groups, specificHandler)
		case *singleLineRunnableHandler:
			commands = append(commands, specificHandler)
		case *dialogContextRunnable:
			dialogs = append(dialogs, specificHandler)
		default:
			panic("Unknown handler type")
		}
	}
	return []*lineGroup{
		listRunnableHandlers(commands, "Commands"),
		listRunnableHandlers(dialogs, "Dialogs"),
		listRunnableHandlers(groups, "Groups"),
	}
}

func listRunnableHandlers(handlers []runnableHandler, groupName string) *lineGroup {
	if len(handlers) == 0 {
		return &lineGroup{}
	}
	var lines []string
	for _, handler := range handlers {
		lines = append(lines, handler.getHelpIndexLine())
	}
	return &lineGroup{
		name:  groupName,
		lines: lines,
	}
}

type handlersForDialog struct {
	handlers []runnableHandler
}

var _ handlers = new(handlersForDialog)

func (hfd *handlersForDialog) getHandlers() []runnableHandler {
	return hfd.handlers
}

func (hfd *handlersForDialog) showHandlers(outputter Outputter) {
	parts := hfd.splitPropertiesAndCommands()
	lgs := lineGroups(parts)
	outputter(lgs.String())
}

func (hfd *handlersForDialog) splitPropertiesAndCommands() []*lineGroup {
	commands := make([]runnableHandler, 0)
	properties := make([]runnableHandler, 0)
	for _, handler := range hfd.handlers {
		switch specificHandler := handler.(type) {
		case *singleLineRunnableHandler:
			commands = append(commands, specificHandler)
		case *dialogPropertyHandler, *dialogListPropertyHandler:
			properties = append(properties, specificHandler)
		default:
			panic("Unknown handler type")
		}
	}
	return []*lineGroup{
		listRunnableHandlers(properties, "Properties that can be set"),
		listRunnableHandlers(commands, "Commands"),
	}
}
