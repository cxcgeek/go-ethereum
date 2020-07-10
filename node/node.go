// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package node

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/internal/debug"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/tsdb/fileutil"
)

// Node is a container on which services can be registered.
type Node struct {
	eventmux          *event.TypeMux // Event multiplexer used between the services of a stack
	config            *Config
	accman            *accounts.Manager
	log               log.Logger
	ephemeralKeystore string            // if non-empty, the key directory that will be removed by Stop
	instanceDirLock   fileutil.Releaser // prevents concurrent use of instance directory

	lock          sync.RWMutex
	stop          chan struct{}              // Channel to wait for termination notifications
	server        *p2p.Server                // Currently running P2P networking layer
	lifecycles    map[reflect.Type]Lifecycle // All registered backends, services, and auxiliary services that have a lifecycle
	httpServers   serverMap                  // serverMap stores information about the node's rpc, ws, and graphQL http servers.
	inprocHandler *rpc.Server                // In-process RPC request handler to process the API requests
	rpcAPIs       []rpc.API                  // List of APIs currently provided by the node
	ipc           *httpServer                // Stores information about the ipc http server
}

// New creates a new P2P node, ready for protocol registration.
func New(conf *Config) (*Node, error) {
	// Copy config and resolve the datadir so future changes to the current
	// working directory don't affect the node.
	confCopy := *conf
	conf = &confCopy
	if conf.DataDir != "" {
		absdatadir, err := filepath.Abs(conf.DataDir)
		if err != nil {
			return nil, err
		}
		conf.DataDir = absdatadir
	}
	// Ensure that the instance name doesn't cause weird conflicts with
	// other files in the data directory.
	if strings.ContainsAny(conf.Name, `/\`) {
		return nil, errors.New(`Config.Name must not contain '/' or '\'`)
	}
	if conf.Name == datadirDefaultKeyStore {
		return nil, errors.New(`Config.Name cannot be "` + datadirDefaultKeyStore + `"`)
	}
	if strings.HasSuffix(conf.Name, ".ipc") {
		return nil, errors.New(`Config.Name cannot end in ".ipc"`)
	}
	// Ensure that the AccountManager method works before the node has started.
	// We rely on this in cmd/geth.
	am, ephemeralKeystore, err := makeAccountManager(conf)
	if err != nil {
		return nil, err
	}
	if conf.Logger == nil {
		conf.Logger = log.New()
	}
	// Note: any interaction with Config that would create/touch files
	// in the data directory or instance directory is delayed until Start.
	node := &Node{
		accman:            am,
		ephemeralKeystore: ephemeralKeystore,
		config:            conf,
		lifecycles:        make(map[reflect.Type]Lifecycle),
		httpServers:       make(serverMap),
		ipc: &httpServer{
			endpoint: conf.IPCEndpoint(),
		},
		inprocHandler: rpc.NewServer(),
		eventmux:      new(event.TypeMux),
		log:           conf.Logger,
	}

	// Initialize the p2p server. This creates the node key and
	// discovery databases.
	node.server = &p2p.Server{Config: conf.P2P}
	node.server.Config.PrivateKey = node.config.NodeKey()
	node.server.Config.Name = node.config.NodeName()
	node.server.Config.Logger = node.log
	if node.server.Config.StaticNodes == nil {
		node.server.Config.StaticNodes = node.config.StaticNodes()
	}
	if node.server.Config.TrustedNodes == nil {
		node.server.Config.TrustedNodes = node.config.TrustedNodes()
	}
	if node.server.Config.NodeDatabase == "" {
		node.server.Config.NodeDatabase = node.config.NodeDB()
	}

	// Configure HTTP server(s)
	if conf.HTTPHost != "" {
		httpServ := &httpServer{
			CorsAllowedOrigins: conf.HTTPCors,
			Vhosts:             conf.HTTPVirtualHosts,
			Whitelist:          conf.HTTPModules,
			Timeouts:           conf.HTTPTimeouts,
			Srv:                rpc.NewServer(),
			endpoint:           conf.HTTPEndpoint(),
			host:               conf.HTTPHost,
			port:               conf.HTTPPort,
			RPCAllowed:         1,
		}
		// check if ws is enabled and if ws port is the same as http port
		if conf.WSHost != "" && conf.WSPort == conf.HTTPPort {
			httpServ.WSAllowed = 1
			httpServ.WsOrigins = conf.WSOrigins
			httpServ.Whitelist = append(httpServ.Whitelist, conf.WSModules...)

			node.httpServers[conf.HTTPEndpoint()] = httpServ
			return node, nil
		}
		node.httpServers[conf.HTTPEndpoint()] = httpServ
	}
	if conf.WSHost != "" {
		node.httpServers[conf.WSEndpoint()] = &httpServer{
			WsOrigins: conf.WSOrigins,
			Whitelist: conf.WSModules,
			Srv:       rpc.NewServer(),
			endpoint:  conf.WSEndpoint(),
			host:      conf.WSHost,
			port:      conf.WSPort,
			WSAllowed: 1,
		}
	}

	return node, nil
}

// Close stops the Node and releases resources acquired in
// Node constructor New.
func (n *Node) Close() error {
	var errs []error

	// Terminate all subsystems and collect any errors
	if err := n.Stop(); err != nil && err != ErrNodeStopped {
		errs = append(errs, err)
	}
	if err := n.accman.Close(); err != nil {
		errs = append(errs, err)
	}
	// Report any errors that might have occurred
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return fmt.Errorf("%v", errs)
	}
}

// RegisterLifecycle registers the given Lifecycle on the node.
func (n *Node) RegisterLifecycle(lifecycle Lifecycle) {
	kind := reflect.TypeOf(lifecycle)
	if _, exists := n.lifecycles[kind]; exists {
		n.Fatalf("Lifecycle cannot be registered more than once", kind)
	}

	n.lifecycles[kind] = lifecycle
}

// RegisterProtocols adds backend's protocols to the node's p2p server.
func (n *Node) RegisterProtocols(protocols []p2p.Protocol) {
	n.server.Protocols = append(n.server.Protocols, protocols...)
}

// RegisterAPIs registers the APIs a service provides on the node.
func (n *Node) RegisterAPIs(apis []rpc.API) {
	n.rpcAPIs = append(n.rpcAPIs, apis...)
}

// RegisterHTTPServer registers the given HTTP server on the node.
func (n *Node) RegisterHTTPServer(endpoint string, server *httpServer) {
	n.httpServers[endpoint] = server
}

// RegisterPath mounts the given handler on the given path on the canonical HTTP server.
func (n *Node) RegisterPath(path string, handler http.Handler) string {
	for _, server := range n.httpServers {
		if atomic.LoadInt32(&server.RPCAllowed) == 1 {
			server.srvMux.Handle(path, handler)
			return server.endpoint
		}
	}
	n.log.Warn(fmt.Sprintf("HTTP server not configured on node, path %s cannot be enabled", path))
	return ""
}

// ExistingHTTPServer checks if an HTTP server is already configured on the given endpoint.
func (n *Node) ExistingHTTPServer(endpoint string) *httpServer {
	if server, exists := n.httpServers[endpoint]; exists {
		return server
	}
	return nil
}

// CreateHTTPServer creates an http.Server and adds it to the given httpServer.
func (n *Node) CreateHTTPServer(h *httpServer, exposeAll bool) error {
	// register apis and create handler stack
	err := RegisterApisFromWhitelist(n.rpcAPIs, h.Whitelist, h.Srv, exposeAll)
	if err != nil {
		return err
	}

	// start the HTTP listener
	listener, err := net.Listen("tcp", h.endpoint)
	if err != nil {
		return err
	}
	// create the HTTP server
	httpSrv := &http.Server{Handler: &h.srvMux}
	// check timeouts if they exist
	if h.Timeouts != (rpc.HTTPTimeouts{}) {
		CheckTimeouts(&h.Timeouts)
		httpSrv.ReadTimeout = h.Timeouts.ReadTimeout
		httpSrv.WriteTimeout = h.Timeouts.WriteTimeout
		httpSrv.IdleTimeout = h.Timeouts.IdleTimeout
	}
	// add listener and http.Server to httpServer
	h.Listener = listener
	h.Server = httpSrv

	return nil
}

// running returns true if the node's p2p server is already running.
func (n *Node) running() bool {
	return n.server.Running()
}

// Start creates a live P2P node and starts running it.
func (n *Node) Start() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	// Short circuit if the node's already running
	if n.running() {
		return ErrNodeRunning
	}
	if err := n.openDataDir(); err != nil {
		return err
	}

	// Start the p2p node
	if err := n.server.Start(); err != nil {
		return convertFileLockError(err)
	}
	n.log.Info("Starting peer-to-peer node", "instance", n.server.Name)

	// Configure the RPC interfaces
	if err := n.configureRPC(); err != nil {
		n.httpServers.Stop()
		n.server.Stop()
		return err
	}

	// Start all registered lifecycles
	var started []Lifecycle
	for _, lifecycle := range n.lifecycles {
		if err := lifecycle.Start(); err != nil {
			n.stopLifecycles(started)
			n.server.Stop()
			return err
		}
		started = append(started, lifecycle)
	}

	// Finish initializing the startup
	n.stop = make(chan struct{})
	return nil
}

// stopLifecycles stops the node's running Lifecycles.
func (n *Node) stopLifecycles(started []Lifecycle) {
	for _, lifecycle := range started {
		lifecycle.Stop()
	}
}

// Config returns the configuration of node.
func (n *Node) Config() *Config {
	return n.config
}

func (n *Node) openDataDir() error {
	if n.config.DataDir == "" {
		return nil // ephemeral
	}

	instdir := filepath.Join(n.config.DataDir, n.config.name())
	if err := os.MkdirAll(instdir, 0700); err != nil {
		return err
	}
	// Lock the instance directory to prevent concurrent use by another instance as well as
	// accidental use of the instance directory as a database.
	release, _, err := fileutil.Flock(filepath.Join(instdir, "LOCK"))
	if err != nil {
		return convertFileLockError(err)
	}
	n.instanceDirLock = release
	return nil
}

// configureRPC is a helper method to configure all the various RPC endpoints during node
// startup. It's not meant to be called at any time afterwards as it makes certain
// assumptions about the state of the node.
func (n *Node) configureRPC() error {
	n.RegisterAPIs(n.apis())

	// Start the various API endpoints, terminating all in case of errors
	if err := n.startInProc(); err != nil {
		return err
	}
	if err := n.startIPC(); err != nil {
		n.stopInProc()
		return err
	}
	// configure HTTPServers
	for _, server := range n.httpServers {
		// configure the handlers
		handler := n.createHandler(server)
		if handler != nil {
			server.srvMux.Handle("/", handler)
		}
		// create the HTTP server
		if err := n.CreateHTTPServer(server, false); err != nil {
			return err
		}
	}
	// only register http server as a lifecycle if it has not already been registered
	if _, exists := n.lifecycles[reflect.TypeOf(n.httpServers)]; !exists {
		n.RegisterLifecycle(n.httpServers)
	}
	// All API endpoints started successfully
	return nil
}

// createHandler creates the http.Handler for the given httpServer.
func (n *Node) createHandler(server *httpServer) http.Handler {
	var handler http.Handler
	if atomic.LoadInt32(&server.RPCAllowed) == 1 {
		handler = NewHTTPHandlerStack(server.Srv, server.CorsAllowedOrigins, server.Vhosts)
		// wrap ws handler just in case ws is enabled through the console after start-up
		wsHandler := server.Srv.WebsocketHandler(server.WsOrigins)
		handler = server.NewWebsocketUpgradeHandler(handler, wsHandler)

		n.log.Info("HTTP configured on endpoint ", "endpoint", fmt.Sprintf("http://%s/", server.endpoint))
		if atomic.LoadInt32(&server.WSAllowed) == 1 {
			n.log.Info("Websocket configured on endpoint ", "endpoint", fmt.Sprintf("ws://%s/", server.endpoint))
		}
	}
	if (atomic.LoadInt32(&server.WSAllowed) == 1) && handler == nil {
		handler = server.Srv.WebsocketHandler(server.WsOrigins)
		n.log.Info("Websocket configured on endpoint ", "endpoint", fmt.Sprintf("ws://%s/", server.endpoint))
	}

	return handler
}

// startInProc registers all RPC APIs on the inproc server.
func (n *Node) startInProc() error {
	for _, api := range n.rpcAPIs {
		if err := n.inprocHandler.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
	}
	return nil
}

// stopInProc terminates the in-process RPC endpoint.
func (n *Node) stopInProc() {
	n.inprocHandler.Stop()
}

// startIPC initializes and starts the IPC RPC endpoint.
func (n *Node) startIPC() error {
	if n.ipc.endpoint == "" {
		return nil // IPC disabled.
	}
	listener, handler, err := rpc.StartIPCEndpoint(n.ipc.endpoint, n.rpcAPIs)
	if err != nil {
		return err
	}
	n.ipc.Listener = listener
	n.ipc.handler = handler
	n.log.Info("IPC endpoint opened", "url", n.ipc.endpoint)
	return nil
}

// stopIPC terminates the IPC RPC endpoint.
func (n *Node) stopIPC() {
	if n.ipc.Listener != nil {
		n.ipc.Listener.Close()
		n.ipc.Listener = nil

		n.log.Info("IPC endpoint closed", "url", n.ipc.endpoint)
	}
	if n.ipc.Srv != nil {
		n.ipc.Srv.Stop()
		n.ipc.Srv = nil
	}
}

// stopServers terminates the given HTTP servers' endpoints
func (n *Node) stopServer(server *httpServer) {
	if server.Server != nil {
		url := fmt.Sprintf("http://%v/", server.Listener.Addr())
		// Don't bother imposing a timeout here.
		server.Server.Shutdown(context.Background())
		n.log.Info("HTTP Endpoint closed", "url", url)
	}
	if server.Srv != nil {
		server.Srv.Stop()
		server.Srv = nil
	}
	// remove stopped http server from node's http servers
	delete(n.httpServers, server.endpoint)
}

// Stop terminates a running node along with all it's services. In the node was
// not started, an error is returned.
func (n *Node) Stop() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	// Short circuit if the node's not running
	if n.server == nil || !n.running() {
		return ErrNodeStopped
	}

	// Terminate the API, services and the p2p server.
	n.stopIPC()
	n.rpcAPIs = nil
	failure := &StopError{
		Services: make(map[reflect.Type]error),
	}
	for kind, lifecycle := range n.lifecycles {
		if err := lifecycle.Stop(); err != nil {
			failure.Services[reflect.TypeOf(lifecycle)] = err
		}
		delete(n.lifecycles, kind)
	}
	n.server.Stop()
	n.server = nil

	// Release instance directory lock.
	if n.instanceDirLock != nil {
		if err := n.instanceDirLock.Release(); err != nil {
			n.log.Error("Can't release datadir lock", "err", err)
		}
		n.instanceDirLock = nil
	}

	// unblock n.Wait
	close(n.stop)

	// Remove the keystore if it was created ephemerally.
	var keystoreErr error
	if n.ephemeralKeystore != "" {
		keystoreErr = os.RemoveAll(n.ephemeralKeystore)
	}

	if len(failure.Services) > 0 {
		return failure
	}
	if keystoreErr != nil {
		return keystoreErr
	}
	return nil
}

// Wait blocks the thread until the node is stopped. If the node is not running
// at the time of invocation, the method immediately returns.
func (n *Node) Wait() {
	n.lock.RLock()
	if n.server == nil {
		n.lock.RUnlock()
		return
	}
	stop := n.stop
	n.lock.RUnlock()

	<-stop
}

// Attach creates an RPC client attached to an in-process API handler.
func (n *Node) Attach() (*rpc.Client, error) {
	return rpc.DialInProc(n.inprocHandler), nil
}

// RPCHandler returns the in-process RPC request handler.
func (n *Node) RPCHandler() (*rpc.Server, error) {
	n.lock.RLock()
	defer n.lock.RUnlock()

	if n.inprocHandler == nil {
		return nil, ErrNodeStopped
	}
	return n.inprocHandler, nil
}

// Server retrieves the currently running P2P network layer. This method is meant
// only to inspect fields of the currently running server, life cycle management
// should be left to this Node entity.
func (n *Node) Server() *p2p.Server {
	n.lock.RLock()
	defer n.lock.RUnlock()

	return n.server
}

// DataDir retrieves the current datadir used by the protocol stack.
// Deprecated: No files should be stored in this directory, use InstanceDir instead.
func (n *Node) DataDir() string {
	return n.config.DataDir
}

// InstanceDir retrieves the instance directory used by the protocol stack.
func (n *Node) InstanceDir() string {
	return n.config.instanceDir()
}

// AccountManager retrieves the account manager used by the protocol stack.
func (n *Node) AccountManager() *accounts.Manager {
	return n.accman
}

// IPCEndpoint retrieves the current IPC endpoint used by the protocol stack.
func (n *Node) IPCEndpoint() string {
	return n.ipc.endpoint
}

// WSEndpoint retrieves the current WS endpoint used by the protocol stack.
func (n *Node) WSEndpoint() string {
	n.lock.Lock()
	defer n.lock.Unlock()

	for _, httpServer := range n.httpServers {
		if atomic.LoadInt32(&httpServer.WSAllowed) == 1 {
			if httpServer.Listener != nil {
				return httpServer.Listener.Addr().String()
			}
			return httpServer.endpoint
		}
	}

	return n.config.WSEndpoint()
}

// EventMux retrieves the event multiplexer used by all the network services in
// the current protocol stack.
func (n *Node) EventMux() *event.TypeMux {
	return n.eventmux
}

// OpenDatabase opens an existing database with the given name (or creates one if no
// previous can be found) from within the node's instance directory. If the node is
// ephemeral, a memory database is returned.
func (n *Node) OpenDatabase(name string, cache, handles int, namespace string) (ethdb.Database, error) {
	if n.config.DataDir == "" {
		return rawdb.NewMemoryDatabase(), nil
	}
	return rawdb.NewLevelDBDatabase(n.config.ResolvePath(name), cache, handles, namespace)
}

// OpenDatabaseWithFreezer opens an existing database with the given name (or
// creates one if no previous can be found) from within the node's data directory,
// also attaching a chain freezer to it that moves ancient chain data from the
// database to immutable append-only files. If the node is an ephemeral one, a
// memory database is returned.
func (n *Node) OpenDatabaseWithFreezer(name string, cache, handles int, freezer, namespace string) (ethdb.Database, error) {
	if n.config.DataDir == "" {
		return rawdb.NewMemoryDatabase(), nil
	}
	root := n.config.ResolvePath(name)

	switch {
	case freezer == "":
		freezer = filepath.Join(root, "ancient")
	case !filepath.IsAbs(freezer):
		freezer = n.config.ResolvePath(freezer)
	}
	return rawdb.NewLevelDBDatabaseWithFreezer(root, cache, handles, freezer, namespace)
}

// ResolvePath returns the absolute path of a resource in the instance directory.
func (n *Node) ResolvePath(x string) string {
	return n.config.ResolvePath(x)
}

// apis returns the collection of RPC descriptors this node offers.
func (n *Node) apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(n),
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPublicAdminAPI(n),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   debug.Handler,
		}, {
			Namespace: "web3",
			Version:   "1.0",
			Service:   NewPublicWeb3API(n),
			Public:    true,
		},
	}
}

// RegisterApisFromWhitelist checks the given modules' availability, generates a whitelist based on the allowed modules,
// and then registers all of the APIs exposed by the services.
func RegisterApisFromWhitelist(apis []rpc.API, modules []string, srv *rpc.Server, exposeAll bool) error {
	if bad, available := checkModuleAvailability(modules, apis); len(bad) > 0 {
		log.Error("Unavailable modules in HTTP API list", "unavailable", bad, "available", available)
	}
	// Generate the whitelist based on the allowed modules
	whitelist := make(map[string]bool)
	for _, module := range modules {
		whitelist[module] = true
	}
	// Register all the APIs exposed by the services
	for _, api := range apis {
		if exposeAll || whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
			if err := srv.RegisterName(api.Namespace, api.Service); err != nil {
				return err
			}
		}
	}
	return nil
}

// Fatalf formats a message to standard error and exits the program.
// The message is also printed to standard output if standard error
// is redirected to a different file.
func (n *Node) Fatalf(format string, args ...interface{}) {
	w := io.MultiWriter(os.Stdout, os.Stderr)
	if runtime.GOOS == "windows" {
		// The SameFile check below doesn't work on Windows.
		// stdout is unlikely to get redirected though, so just print there.
		w = os.Stdout
	} else {
		outf, _ := os.Stdout.Stat()
		errf, _ := os.Stderr.Stat()
		if outf != nil && errf != nil && os.SameFile(outf, errf) {
			w = os.Stderr
		}
	}
	fmt.Fprintf(w, "Fatal: "+format+"\n", args...)
	os.Exit(1)
}
