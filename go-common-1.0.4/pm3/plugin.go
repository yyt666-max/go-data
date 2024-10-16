package pm3

type tEmptyDriver struct {
	plugin IPlugin
}

func EmptyDriver(plugin IPlugin) Driver {
	return &tEmptyDriver{plugin: plugin}
}

func (e *tEmptyDriver) Create() (IPlugin, error) {
	return e.plugin, nil
}

func (e *tEmptyDriver) Access() map[string][]string {
	return nil
}
