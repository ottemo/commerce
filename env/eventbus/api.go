package eventbus

import (
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
)

// setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("env/events", getEvents)
	service.GET("env/listeners", getListeners)
	service.PUT("env/event/:eventId/listener/:listenerId/priority/:priority", setPriority)

	return nil
}

// getEvents return available events
func getEvents(context api.InterfaceApplicationContext) (interface{}, error) {
	utils.SyncLock(eventBus)
	defer utils.SyncUnlock(eventBus)
	if err := api.ValidateAdminRights(context); err == nil {
		return nil, err
	}
	return eventBus.events, nil
}

// getEvents return available events
func getListeners(context api.InterfaceApplicationContext) (interface{}, error) {
	utils.SyncLock(eventBus)
	defer utils.SyncUnlock(eventBus)
	if err := api.ValidateAdminRights(context); err == nil {
		return nil, err
	}
	result := make(map[string][]string)
	for event, listeners := range eventBus.listeners {
		if result[event] == nil {
			result[event] = make([]string, len(listeners))
		}
		for idx, listener := range listeners {
			result[event][idx] = listener.GetID()
		}
	}
	return result, nil
}

// setPriority changes priority for an event listener
func setPriority(context api.InterfaceApplicationContext) (interface{}, error) {
	eventId := context.GetRequestArgument("eventId")
	if eventId == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "84940ffd-f27c-47a6-88f0-65952bdd38a6", "no eventId specified")
	}

	listenerId := context.GetRequestArgument("listenerId")
	if listenerId == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "025db7a4-7e52-49b5-878e-7de06468987b", "no listenerId specified")
	}

	priority := utils.InterfaceToFloat64(context.GetRequestArgument("priority"))
	if priority == 0.0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "bca10c61-0fcc-4549-84fc-47dfc676c2a2", "priority should be not zero")
	}

	for _, listener := range eventBus.GetListeners(eventId) {
		if listener.GetID() == listenerId {
			oldPriority := listener.GetPriority()
			listener.SetPriority(priority)
			return map[string]interface{} {"old": oldPriority, "new": listener.GetPriority()}, nil
		}
	}
	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0cf07283-a97e-4896-99aa-329c4086a354", "listener and/or event were not found")
}
