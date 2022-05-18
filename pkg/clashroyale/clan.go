package clashroyale

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	goclash "github.com/fiskWasTaken/go-clash"

	"github.com/mecode4food/cr-clan-bot/pkg/config"
)

var (
	c *goclash.Client
	m = sync.RWMutex{}

	token = config.Viper().GetString("clash_royale.token")
)

func Client() (*goclash.Client, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	m.RLock()
	if c != nil {
		defer m.RUnlock()
		return c, nil
	}
	m.RUnlock()

	m.Lock()
	defer m.Unlock()
	if c != nil {
		defer m.Unlock()
		return c, nil
	}

	var err error
	c = goclash.NewClient(token)
	base, err := url.Parse("https://proxy.royaleapi.dev")
	if err != nil {
		return nil, err
	}
	c.BaseURL = base

	return c, nil
}

func Clan(id string) (goclash.Clan, error) {
	cl, err := Client()
	if err != nil {
		return goclash.Clan{}, err
	}

	cs := cl.Clan(id)
	c, err := cs.Get()
	if err != nil {
		return goclash.Clan{}, err
	}

	return c, nil
}

func ClanMembers(id string) ([]goclash.ClanMember, error) {
	cl, err := Client()
	if err != nil {
		return nil, err
	}

	cs := cl.Clan(id)

	// var mm []goclash.ClanMember
	mp, err := cs.Members()
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v\n", mp)
	return mp.Items, nil
}
