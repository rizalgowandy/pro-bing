//go:build linux
// +build linux

package probing

import (
	"errors"
	"os"
	"reflect"
	"syscall"

	"golang.org/x/net/icmp"
)

// Returns the length of an ICMP message.
func (p *Pinger) getMessageLength() int {
	return p.Size + 8
}

// Attempts to match the ID of an ICMP packet.
func (p *Pinger) matchID(ID int) bool {
	// On Linux we can only match ID if we are privileged.
	if p.protocol == "icmp" {
		return ID == p.id
	}
	return true
}

// SetMark sets the SO_MARK socket option on outgoing ICMP packets.
// Setting this option requires CAP_NET_ADMIN.
func (c *icmpConn) SetMark(mark uint) error {
	fd, err := getFD(c.c)
	if err != nil {
		return err
	}
	return os.NewSyscallError(
		"setsockopt",
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(mark)),
	)
}

// SetMark sets the SO_MARK socket option on outgoing ICMP packets.
// Setting this option requires CAP_NET_ADMIN.
func (c *icmpv4Conn) SetMark(mark uint) error {
	fd, err := getFD(c.icmpConn.c)
	if err != nil {
		return err
	}
	return os.NewSyscallError(
		"setsockopt",
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(mark)),
	)
}

// SetMark sets the SO_MARK socket option on outgoing ICMP packets.
// Setting this option requires CAP_NET_ADMIN.
func (c *icmpV6Conn) SetMark(mark uint) error {
	fd, err := getFD(c.icmpConn.c)
	if err != nil {
		return err
	}
	return os.NewSyscallError(
		"setsockopt",
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(mark)),
	)
}

// getFD gets the system file descriptor for an icmp.PacketConn
func getFD(c *icmp.PacketConn) (uintptr, error) {
	v := reflect.ValueOf(c).Elem().FieldByName("c").Elem()
	if v.Elem().Kind() != reflect.Struct {
		return 0, errors.New("invalid type")
	}

	fd := v.Elem().FieldByName("conn").FieldByName("fd")
	if fd.Elem().Kind() != reflect.Struct {
		return 0, errors.New("invalid type")
	}

	pfd := fd.Elem().FieldByName("pfd")
	if pfd.Kind() != reflect.Struct {
		return 0, errors.New("invalid type")
	}

	return uintptr(pfd.FieldByName("Sysfd").Int()), nil
}
