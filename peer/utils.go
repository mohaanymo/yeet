package peer

import (
	"os/user"
	"net"
)

func GetUserName() (string) {
	currentUser, err := user.Current()
	if err!=nil{
		return "Uknown User"
	}
	return currentUser.Username
}

func GetLocalIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "localhost"
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP.String()
}