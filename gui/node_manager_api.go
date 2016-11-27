package gui

import
//"encoding/json"

//"strconv"

(
	"encoding/json"
	"io"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/cxo/nodeManager"
)

//RegisterNodeManagerHandlers - create routes for NodeManager
func RegisterNodeManagerHandlers(router *Router, shm *nodeManager.Manager) {
	// enclose shm into SkyhashManager to be able to add methods
	lshm := SkyhashManager{Manager: shm}

	router.GET("/manager", lshm._ManagerInfo)
	router.POST("/manager", lshm._StartManager)
	router.DELETE("/manager", lshm._StopManager)

	router.GET("/manager/nodes", lshm._ListNodes)
	router.POST("/manager/nodes", lshm._AddNode)

	router.GET("/manager/nodes/:node_id", lshm._GetNodeByID)
	router.DELETE("/manager/nodes/:node_id", lshm._StopNodeByID)

	router.GET("/manager/nodes/:node_id/subscriptions", lshm._ListSubscriptions)
	router.POST("/manager/nodes/:node_id/subscriptions", lshm._AddSubscription)

	router.GET("/manager/nodes/:node_id/subscriptions/:subscription_id", lshm._GetSubscriptionByID)
	router.DELETE("/manager/nodes/:node_id/subscriptions/:subscription_id", lshm._RemoveSubscriptionByID)

	router.GET("/manager/nodes/:node_id/subscribers", lshm._ListSubscribers)

	router.GET("/manager/nodes/:node_id/subscribers/:subscriber_id", lshm._GetSubscriberByID)
	router.DELETE("/manager/nodes/:node_id/subscribers/:subscriber_id", lshm._RemoveSubscriberByID)

}

func (self *SkyhashManager) _ManagerInfo(ctx *Context) error {
	return ctx.Text(200, "Hello world!")
}
func (self *SkyhashManager) _StartManager(ctx *Context) error {
	return ctx.ErrInternal("NOT IMPLEMENTED, YET")
}
func (self *SkyhashManager) _StopManager(ctx *Context) error {
	return ctx.ErrInternal("NOT IMPLEMENTED, YET")
}

func (self *SkyhashManager) _ListNodes(ctx *Context) error {
	listOfNodes := self.Nodes()
	OutputNodes := []nodeManager.NodeJSON{}
	for _, node := range listOfNodes {
		OutputNodes = append(OutputNodes, node.JSON())
	}
	return ctx.JSON(200, OutputNodes)
}

func (self *SkyhashManager) _AddNode(ctx *Context) error {

	// decode configuration of the new transport to be created
	var newNodeConfig struct {
		SecKey string `json:"secKey"`
	}

	var newNode *nodeManager.Node
	err := json.NewDecoder(ctx.Request.Body).Decode(&newNodeConfig)
	if err != nil {
		if err == io.EOF {
			newNode = self.NewNode()
		} else {
			return ctx.ErrInternal(err.Error())
		}
	} else {

		secKey, err := cipher.SecKeyFromHex(newNodeConfig.SecKey)
		if err != nil {
			return ctx.ErrInvalidRequest(err.Error(), "secKey", newNodeConfig.SecKey)
		}

		if secKey == (cipher.SecKey{}) {
			return ctx.ErrInvalidRequest("invalid secKey", "secKey", secKey.Hex())
		}

		newNode, err = self.NewNodeFromSecKey(secKey)
		if err != nil {
			return ctx.ErrInvalidRequest("error creating new node from key pair", "err", err)
		}
	}

	err = self.AddNode(newNode)
	if err != nil {
		return ctx.ErrInternal(err.Error())
	}

	// TODO: add callback to register messages, DB callback, etc.
	// otherwise this node has no messages registered

	err = newNode.Start()
	if err != nil {
		return ctx.ErrInternal(err.Error())
	}

	response := JSONResponse{
		Code:   "created",
		Status: 200,
		Detail: "The node has been created",
	}
	return ctx.JSON(200, response)
}

func (self *SkyhashManager) _GetNodeByID(ctx *Context) error {
	nodePubKey, err := ctx.PubKeyFromParam("node_id")
	if err != nil {
		return ctx.ErrInvalidRequest(err.Error())
	}

	node, err := self.NodeByID(*nodePubKey)
	if err != nil {
		return ctx.ErrNotFound(err.Error(), "pubKey", nodePubKey.Hex())
	}

	return ctx.JSON(200, node.JSON())
}

func (self *SkyhashManager) _StopNodeByID(ctx *Context) error {
	return ctx.ErrInternal("NOT IMPLEMENTED, YET")
}

func (self *SkyhashManager) _ListSubscriptions(ctx *Context) error {
	nodePubKey, err := ctx.PubKeyFromParam("node_id")
	if err != nil {
		return ctx.ErrInvalidRequest(err.Error())
	}

	node, err := self.NodeByID(*nodePubKey)
	if err != nil {
		return ctx.ErrNotFound(err.Error(), "pubKey", nodePubKey.Hex())
	}
	listOfSubscriptions := node.Subscriptions()

	OutputSubscriptions := []nodeManager.SubscriptionJSON{}

	for _, subscription := range listOfSubscriptions {
		OutputSubscriptions = append(OutputSubscriptions, subscription.JSON())
	}

	return ctx.JSON(200, OutputSubscriptions)
}

func (self *SkyhashManager) _AddSubscription(ctx *Context) error {
	return ctx.ErrInternal("NOT IMPLEMENTED, YET")
}

func (self *SkyhashManager) _GetSubscriptionByID(ctx *Context) error {
	nodePubKey, err := ctx.PubKeyFromParam("node_id")
	if err != nil {
		return ctx.ErrInvalidRequest(err.Error())
	}

	node, err := self.NodeByID(*nodePubKey)
	if err != nil {
		return ctx.ErrNotFound(err.Error(), "pubKey", nodePubKey.Hex())
	}

	subscriptionPubKey, err := ctx.PubKeyFromParam("subscription_id")
	if err != nil {
		return ctx.ErrInvalidRequest(err.Error())
	}

	subscription, err := node.SubscriptionByID(subscriptionPubKey)
	if err != nil {
		return ctx.ErrNotFound(err.Error(), "pubKey", subscriptionPubKey.Hex())
	}

	return ctx.JSON(200, subscription.JSON())
}

func (self *SkyhashManager) _RemoveSubscriptionByID(ctx *Context) error {
	return ctx.ErrInternal("NOT IMPLEMENTED, YET")
}

func (self *SkyhashManager) _ListSubscribers(ctx *Context) error {
	nodePubKey, err := ctx.PubKeyFromParam("node_id")
	if err != nil {
		return ctx.ErrInvalidRequest(err.Error())
	}

	node, err := self.NodeByID(*nodePubKey)
	if err != nil {
		return ctx.ErrNotFound(err.Error(), "pubKey", nodePubKey.Hex())
	}
	listOfSubscribers := node.Subscribers()

	OutputSubscribers := []nodeManager.SubscriberJSON{}

	for _, subscriber := range listOfSubscribers {
		OutputSubscribers = append(OutputSubscribers, subscriber.JSON())
	}

	return ctx.JSON(200, OutputSubscribers)
}

func (self *SkyhashManager) _GetSubscriberByID(ctx *Context) error {
	nodePubKey, err := ctx.PubKeyFromParam("node_id")
	if err != nil {
		return ctx.ErrInvalidRequest(err.Error())
	}

	node, err := self.NodeByID(*nodePubKey)
	if err != nil {
		return ctx.ErrNotFound(err.Error(), "pubKey", nodePubKey.Hex())
	}

	subscriberPubKey, err := ctx.PubKeyFromParam("subscriber_id")
	if err != nil {
		return ctx.ErrInvalidRequest(err.Error())
	}

	subscriber, err := node.SubscriberByID(subscriberPubKey)
	if err != nil {
		return ctx.ErrNotFound(err.Error(), "pubKey", subscriberPubKey.Hex())
	}

	return ctx.JSON(200, subscriber.JSON())
}

func (self *SkyhashManager) _RemoveSubscriberByID(ctx *Context) error {
	return ctx.ErrInternal("NOT IMPLEMENTED, YET")
}
