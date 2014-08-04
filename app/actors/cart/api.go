package cart

import(
	"errors"

	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/utils"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/visitor"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("cart", "GET", "info", restCartInfo)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cart", "POST", "add/:productId/:qty", restCartAdd)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cart", "PUT", "update/:itemIdx/:qty", restCartUpdate)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cart", "DELETE", "delete/:itemIdx", restCartDelete)
	if err != nil {
		return err
	}


	return nil
}



// returns cart for current session
func getCurrentCart(params *api.T_APIHandlerParams) (cart.I_Cart, error) {
	sessionCartId := params.Session.Get( cart.SESSION_KEY_CURRENT_CART )

	if sessionCartId != nil {
		currentCart, err := cart.LoadCartById( utils.InterfaceToString(sessionCartId) )
		if err != nil {
			return nil, err
		}

		return currentCart, nil
	} else {

		visitorId := params.Session.Get( visitor.SESSION_KEY_VISITOR_ID )
		if visitorId != nil {
			currentCart, err := cart.GetCartForVisitor( utils.InterfaceToString(visitorId) )
			if err != nil {
				return nil, err
			}

			return currentCart, nil
		} else {
			return nil, errors.New("you are not registered")
		}

	}

	return nil, errors.New("can't get cart for current session")
}



// WEB REST API function to get cart information
//   - parent categories and categorys will not be present in list
func restCartInfo(params *api.T_APIHandlerParams) (interface{}, error) {

	currentCart, err  := getCurrentCart(params)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{} {
		"visitor_id": currentCart.GetVisitorId(),
		"cart_info": currentCart.GetCartInfo(),
		"items": currentCart.ListItems(),
	}

	return result, nil
}



// WEB REST API for adding new item into cart
//   - "pid" (product id) should be specified
//   - "qty" and "options" are optional params
func restCartAdd(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	var pid string = ""
	reqPid, present := reqData["pid"]
	pid = utils.InterfaceToString(reqPid)
	if !present || pid == "" {
		return nil, errors.New("pid should be specified")
	}

	var qty int = 1
	reqQty, present := reqData["qty"]
	if present {
		qty = utils.InterfaceToInt(reqQty)
	}

	var options map[string]interface{} = nil
	reqOptions, present := reqData["options"]
	if present {
		if tmpOptions, ok := reqOptions.(map[string]interface{}); ok {
			options = tmpOptions
		}
	}

	// operation
	//----------
	currentCart, err  := getCurrentCart(params)
	if err != nil {
		return nil, err
	}

	currentCart.AddItem(pid, qty, options)
	currentCart.Save()

	return "ok", nil
}



// WEB REST API used to update cart item qty
//   - "itemIdx" and "qty" should be specified in request URI
func restCartUpdate(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	if !utils.KeysInMapAndNotBlank(params.RequestURLParams, "itemIdx", "qty") {
		return nil, errors.New("itemIdx and qty should be specified")
	}

	itemIdx, err :=  utils.StrToInt(params.RequestURLParams["itemIdx"])
	if err != nil {
		return nil, err
	}

	qty, err :=  utils.StrToInt(params.RequestURLParams["qty"])
	if err != nil {
		return nil, err
	}
	if qty <= 0 {
		return nil, errors.New("qty should be greather then 0")
	}

	// operation
	//----------
	currentCart, err := getCurrentCart(params)
	if err != nil {
		return nil, err
	}

	cartItems := currentCart.ListItems()
	if len(cartItems) > itemIdx {
		cartItems[itemIdx].SetQty( qty )
	} else {
		return nil, errors.New("wrong itemIdx was specified")
	}
	currentCart.Save()

	return "ok", nil
}



// WEB REST API used to delete cart item from cart
//   - "itemIdx" should be specified in request URI
func restCartDelete(params *api.T_APIHandlerParams) (interface{}, error) {

	_, present := params.RequestURLParams["itemIdx"]
	if !present {
		return nil, errors.New("itemIdx should be specified")
	}

	itemIdx, err :=  utils.StrToInt(params.RequestURLParams["itemIdx"])
	if err != nil {
		return nil, err
	}

	// operation
	//----------
	currentCart, err := getCurrentCart(params)
	if err != nil {
		return nil, err
	}

	err = currentCart.RemoveItem(itemIdx)
	if err != nil {
		return nil, err
	}
	currentCart.Save()

	return "ok", nil
}
