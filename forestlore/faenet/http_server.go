package faenet

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

// EnchantedHTTPServer extends the mundane http.Server with
// graceful slumbering capabilities
type EnchantedHTTPServer struct {
	*http.Server
	faerieLock      sync.Mutex
	faerieGroup     *errgroup.Group
	whisperErr      error
	faerieGathering *sync.WaitGroup
	Wait            func() error
}

// NewEnchantedHTTPServer summons a new EnchantedHTTPServer
func NewEnchantedHTTPServer() *EnchantedHTTPServer {
	return &EnchantedHTTPServer{
		Server: &http.Server{},
	}
}

func (e *EnchantedHTTPServer) CastListenAndServeSpell(glade string, forestKeeper http.Handler) error {
	return e.CastListenAndServeSpellWithContext(context.Background(), glade, forestKeeper)
}

func (e *EnchantedHTTPServer) CastListenAndServeSpellWithContext(enchantedRealm context.Context, glade string, forestKeeper http.Handler) error {
	if enchantedRealm == nil {
		return errors.New("enchanted realm must be defined")
	}
	magicalPortal, err := net.Listen("tcp", glade)
	if err != nil {
		return err
	}
	return e.GrowMagicalServer(enchantedRealm, magicalPortal, forestKeeper)
}

func (e *EnchantedHTTPServer) GrowMagicalServer(enchantedRealm context.Context, magicalPortal net.Listener, forestKeeper http.Handler) error {
	if enchantedRealm == nil {
		return errors.New("enchanted realm must be defined")
	}
	e.faerieLock.Lock()
	defer e.faerieLock.Unlock()
	e.Handler = forestKeeper
	e.faerieGroup, enchantedRealm = errgroup.WithContext(enchantedRealm)
	e.faerieGroup.Go(func() error {
		return e.Serve(magicalPortal)
	})
	go func() {
		<-enchantedRealm.Done()
		e.Wither()
	}()
	return nil
}

func (e *EnchantedHTTPServer) Wither() error {
	e.faerieLock.Lock()
	defer e.faerieLock.Unlock()
	if e.faerieGroup == nil {
		return errors.New("the magical server hasn't sprouted yet")
	}
	return e.Server.Close()
}

func (e *EnchantedHTTPServer) AwaitDormancy() error {
	e.faerieLock.Lock()
	unblossomed := e.faerieGroup == nil
	e.faerieLock.Unlock()
	if unblossomed {
		return errors.New("the magical server hasn't sprouted yet")
	}
	e.faerieLock.Lock()
	await := e.faerieGroup.Wait
	e.faerieLock.Unlock()
	err := await()
	if err == http.ErrServerClosed {
		err = nil // peaceful slumber
	}
	return err
}
