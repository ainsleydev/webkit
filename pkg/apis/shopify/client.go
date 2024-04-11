package shopify

import (
	"encoding/json"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"

	"github.com/ainsleydev/decspets/types"
)

type Client struct {
	shopify *goshopify.Client
}

func New(cfg *types.Config) *Client {
	app := goshopify.App{
		ApiKey:      cfg.ShopifyAPIKey,
		ApiSecret:   cfg.ShopifyAccessToken,
		RedirectUrl: "https://example.com/shopify/callback",
		Scope:       "read_products,read_orders,read_content",
	}
	return &Client{
		shopify: goshopify.NewClient(app, "ainsley-test-store", cfg.ShopifyAccessToken),
	}
}

func (c Client) DoSomething() error {
	//list, err := c.shopify.Product.List(goshopify.ListOptions{})
	//if err != nil {
	//	return err
	//}
	//
	//for _, v := range list {
	//	buf, err := json.MarshalIndent(v, "", "    ")
	//	if err != nil {
	//		log.Error(err)
	//	}
	//	fmt.Println(string(buf))
	//}

	get, err := c.shopify.Page.List(nil)
	if err != nil {
		return err
	}
	log(get)

	vv, _ := c.shopify.Page.ListMetafields(90594279485, nil)
	log(vv)

	return nil
}

func log(response any) {
	buf, _ := json.MarshalIndent(response, "", "    ")
	fmt.Println(string(buf))

}
