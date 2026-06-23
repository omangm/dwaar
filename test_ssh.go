package main
import (
	"fmt"
	"github.com/omangm/dwaar/internal/tunnel"
	"time"
	"net"
)
func main() {
	mgr := tunnel.NewManager()
	rule := tunnel.ForwardRule{
		ID:         "test",
		Name:       "test",
		LocalPort:  8010,
		RemoteHost: "localhost",
		RemotePort: 8010,
		Via: &tunnel.SSHHost{
			Host:         "pop-os",
			Port:         22,
			User:         "omangm",
			IdentityFile: "/home/omangm/.ssh/id_ed25519",
		},
	}
	mgr.Start(rule)
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8010")
	if err != nil {
		fmt.Println("Dial err:", err)
	} else {
		fmt.Println("Dial success")
		conn.Close()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("Logs:", mgr.GetLogs("test"))
}
