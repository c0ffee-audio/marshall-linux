package ble

import (
	"fmt"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

const (
	bluezService   = "org.bluez"
	bluezInterface = "org.bluez.GattCharacteristic1"
	objManager     = "org.freedesktop.DBus.ObjectManager"
	agentPath      = dbus.ObjectPath("/marshall/agent")
)

// agent implémente org.bluez.Agent1 en mode NoInputNoOutput (Just Works).
type agent struct{}

func (a *agent) Release() *dbus.Error                          { return nil }
func (a *agent) RequestPinCode(_ dbus.ObjectPath) (string, *dbus.Error) {
	return "0000", nil
}
func (a *agent) DisplayPinCode(_ dbus.ObjectPath, _ string) *dbus.Error { return nil }
func (a *agent) RequestPasskey(_ dbus.ObjectPath) (uint32, *dbus.Error) {
	return 0, nil
}
func (a *agent) DisplayPasskey(_ dbus.ObjectPath, _ uint32, _ uint16) *dbus.Error { return nil }
func (a *agent) RequestConfirmation(_ dbus.ObjectPath, _ uint32) *dbus.Error      { return nil }
func (a *agent) RequestAuthorization(_ dbus.ObjectPath) *dbus.Error               { return nil }
func (a *agent) AuthorizeService(_ dbus.ObjectPath, _ string) *dbus.Error         { return nil }
func (a *agent) Cancel() *dbus.Error                                              { return nil }

const agentIntrospect = `
<node>
  <interface name="org.bluez.Agent1">
    <method name="Release"/>
    <method name="RequestPinCode"><arg direction="in" type="o"/><arg direction="out" type="s"/></method>
    <method name="DisplayPinCode"><arg direction="in" type="o"/><arg direction="in" type="s"/></method>
    <method name="RequestPasskey"><arg direction="in" type="o"/><arg direction="out" type="u"/></method>
    <method name="DisplayPasskey"><arg direction="in" type="o"/><arg direction="in" type="u"/><arg direction="in" type="q"/></method>
    <method name="RequestConfirmation"><arg direction="in" type="o"/><arg direction="in" type="u"/></method>
    <method name="RequestAuthorization"><arg direction="in" type="o"/></method>
    <method name="AuthorizeService"><arg direction="in" type="o"/><arg direction="in" type="s"/></method>
    <method name="Cancel"/>
  </interface>` + introspect.IntrospectDataString + `</node>`

func registerAgent(conn *dbus.Conn) error {
	conn.Export(&agent{}, agentPath, "org.bluez.Agent1")
	conn.Export(introspect.Introspectable(agentIntrospect), agentPath, "org.freedesktop.DBus.Introspectable")
	mgr := conn.Object(bluezService, "/org/bluez")
	if err := mgr.Call("org.bluez.AgentManager1.RegisterAgent", 0, agentPath, "NoInputNoOutput").Err; err != nil {
		return fmt.Errorf("RegisterAgent: %w", err)
	}
	mgr.Call("org.bluez.AgentManager1.RequestDefaultAgent", 0, agentPath)
	return nil
}

type Client struct {
	conn    *dbus.Conn
	devPath string
	chars   map[string]dbus.ObjectPath
}

// Connect prend soit une adresse MAC soit un nom d'appareil (ex: "MOTIF II A.N.C.").
func Connect(target string) (*Client, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("dbus system bus: %w", err)
	}

	if err := registerAgent(conn); err != nil {
		fmt.Printf("warning: could not register BLE agent: %v\n", err)
	}

	devPath, err := resolveDevice(conn, target)
	if err != nil {
		return nil, err
	}

	if err := ensureConnected(conn, devPath); err != nil {
		return nil, fmt.Errorf("connect BLE: %w", err)
	}

	if err := ensurePaired(conn, devPath); err != nil {
		fmt.Printf("warning: pairing failed: %v\n", err)
	}

	chars, err := discoverCharacteristics(conn, devPath)
	if err != nil {
		return nil, fmt.Errorf("discover characteristics: %w", err)
	}
	if len(chars) == 0 {
		return nil, fmt.Errorf("no GATT characteristics found - is the device paired?")
	}

	return &Client{conn: conn, devPath: devPath, chars: chars}, nil
}

// resolveDevice trouve le path D-Bus de l'appareil par adresse MAC ou par nom.
// Si l'appareil n'est pas en cache, lance un scan BLE court.
func resolveDevice(conn *dbus.Conn, target string) (string, error) {
	isMac := strings.Count(target, ":") == 5

	// Toujours scanner pour trouver le device BLE (adresse aléatoire)
	// car l'adresse BLE peut changer entre les sessions
	fmt.Printf("scanning for %q...\n", target)
	adapter := conn.Object(bluezService, "/org/bluez/hci0")
	filter := map[string]dbus.Variant{
		"Transport": dbus.MakeVariant("le"),
	}
	adapter.Call("org.bluez.Adapter1.SetDiscoveryFilter", 0, filter)
	adapter.Call("org.bluez.Adapter1.StartDiscovery", 0)

	// scan jusqu'à 8s ou jusqu'à trouver le device BLE
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(500 * time.Millisecond)
		if path := findBLEDevice(conn, target, isMac); path != "" {
			adapter.Call("org.bluez.Adapter1.StopDiscovery", 0)
			fmt.Printf("found at %s\n", path)
			return path, nil
		}
	}
	adapter.Call("org.bluez.Adapter1.StopDiscovery", 0)

	return "", fmt.Errorf("device %q not found - make sure it's on and in range", target)
}

// findBLEDevice cherche uniquement les appareils BLE (AddressType = random).
func findBLEDevice(conn *dbus.Conn, target string, isMac bool) string {
	obj := conn.Object(bluezService, "/")
	result := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	if err := obj.Call(objManager+".GetManagedObjects", 0).Store(&result); err != nil {
		return ""
	}

	for path, ifaces := range result {
		devIface, ok := ifaces["org.bluez.Device1"]
		if !ok {
			continue
		}
		addrType, _ := devIface["AddressType"].Value().(string)
		if addrType != "random" {
			continue // on veut seulement le BLE
		}
		if isMac {
			addr, _ := devIface["Address"].Value().(string)
			if strings.EqualFold(addr, target) {
				return string(path)
			}
		} else {
			name, _ := devIface["Name"].Value().(string)
			if strings.EqualFold(name, target) ||
				strings.Contains(strings.ToLower(name), strings.ToLower(target)) {
				return string(path)
			}
		}
	}
	return ""
}

func findInObjects(conn *dbus.Conn, target string, isMac bool) string { //nolint:unused
	obj := conn.Object(bluezService, "/")
	result := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	if err := obj.Call(objManager+".GetManagedObjects", 0).Store(&result); err != nil {
		return ""
	}

	var found string
	for path, ifaces := range result {
		devIface, ok := ifaces["org.bluez.Device1"]
		if !ok {
			continue
		}
		var match bool
		if isMac {
			addr, _ := devIface["Address"].Value().(string)
			match = strings.EqualFold(addr, target)
		} else {
			name, _ := devIface["Name"].Value().(string)
			match = strings.EqualFold(name, target) ||
				strings.Contains(strings.ToLower(name), strings.ToLower(target))
		}
		if !match {
			continue
		}
		// preferrer les appareils BLE (AddressType random ou public BLE)
		// les appareils BR/EDR n'ont pas de caractéristiques GATT dans BlueZ
		addrType, _ := devIface["AddressType"].Value().(string)
		if addrType == "random" {
			return string(path) // BLE aléatoire = priorité max
		}
		found = string(path) // garder en fallback
	}
	return found
}

func ensurePaired(conn *dbus.Conn, devPath string) error {
	obj := conn.Object(bluezService, dbus.ObjectPath(devPath))

	var bonded bool
	obj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.bluez.Device1", "Bonded").Store(&bonded)
	if bonded {
		return nil
	}

	fmt.Println("pairing BLE device...")
	ch := make(chan error, 1)
	go func() {
		ch <- obj.Call("org.bluez.Device1.Pair", 0).Err
	}()
	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("Pair: %w", err)
		}
		fmt.Println("BLE paired successfully")
		time.Sleep(500 * time.Millisecond)
		return nil
	case <-time.After(15 * time.Second):
		return fmt.Errorf("Pair: timeout")
	}
}

func ensureConnected(conn *dbus.Conn, devPath string) error {
	obj := conn.Object(bluezService, dbus.ObjectPath(devPath))

	var connected bool
	if err := obj.Call("org.freedesktop.DBus.Properties.Get", 0,
		"org.bluez.Device1", "Connected").Store(&connected); err != nil {
		return fmt.Errorf("get Connected: %w", err)
	}

	if !connected {
		fmt.Println("connecting...")
		if err := obj.Call("org.bluez.Device1.Connect", 0).Err; err != nil {
			return fmt.Errorf("Connect: %w", err)
		}
	}

	// attendre ServicesResolved (max 10s)
	for i := 0; i < 20; i++ {
		var resolved bool
		obj.Call("org.freedesktop.DBus.Properties.Get", 0,
			"org.bluez.Device1", "ServicesResolved").Store(&resolved)
		if resolved {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for GATT services")
}

func discoverCharacteristics(conn *dbus.Conn, devPath string) (map[string]dbus.ObjectPath, error) {
	obj := conn.Object(bluezService, "/")
	result := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)

	if err := obj.Call(objManager+".GetManagedObjects", 0).Store(&result); err != nil {
		return nil, fmt.Errorf("GetManagedObjects: %w", err)
	}

	chars := make(map[string]dbus.ObjectPath)
	for path, ifaces := range result {
		if !strings.HasPrefix(string(path), devPath) {
			continue
		}
		charIface, ok := ifaces[bluezInterface]
		if !ok {
			continue
		}
		uuid, _ := charIface["UUID"].Value().(string)
		if uuid != "" {
			chars[strings.ToLower(uuid)] = path
		}
	}
	return chars, nil
}

func (c *Client) Read(uuid string) ([]byte, error) {
	path, ok := c.chars[strings.ToLower(uuid)]
	if !ok {
		return nil, fmt.Errorf("characteristic %s not found", uuid)
	}
	obj := c.conn.Object(bluezService, path)

	type readResult struct {
		data []byte
		err  error
	}
	ch := make(chan readResult, 1)
	go func() {
		var result []byte
		err := obj.Call(bluezInterface+".ReadValue", 0, map[string]dbus.Variant{}).Store(&result)
		if err != nil {
			ch <- readResult{nil, fmt.Errorf("ReadValue %s: %w", uuid, err)}
			return
		}
		ch <- readResult{result, nil}
	}()

	select {
	case r := <-ch:
		return r.data, r.err
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("ReadValue %s: timeout", uuid)
	}
}

func (c *Client) Write(uuid string, data []byte) error {
	path, ok := c.chars[strings.ToLower(uuid)]
	if !ok {
		return fmt.Errorf("characteristic %s not found", uuid)
	}
	obj := c.conn.Object(bluezService, path)

	ch := make(chan error, 1)
	go func() {
		ch <- obj.Call(bluezInterface+".WriteValue", 0, data, map[string]dbus.Variant{}).Err
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("WriteValue %s: timeout", uuid)
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) ListCharacteristics() []string {
	uuids := make([]string, 0, len(c.chars))
	for uuid := range c.chars {
		uuids = append(uuids, uuid)
	}
	return uuids
}
