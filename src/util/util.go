package util

import "net"

func GetLocalIPAddress() (string, error) {
	var localAddress string
	conn, err := net.Dial("udp", "1.1.1.1:80")

	if err != nil {
		return localAddress, err
	}

	localAddress = conn.LocalAddr().(*net.UDPAddr).IP.String()

	return localAddress, nil
}
