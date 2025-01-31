// Package router_address implements the I2P RouterAddress common data structure
package router_address

import (
	"errors"

	. "github.com/go-i2p/go-i2p/lib/common/data"
	log "github.com/sirupsen/logrus"
)

// Minimum number of bytes in a valid RouterAddress
const (
	ROUTER_ADDRESS_MIN_SIZE = 9
)

/*
[RouterAddress]
Accurate for version 0.9.49

Description
This structure defines the means to contact a router through a transport protocol.

Contents
1 byte Integer defining the relative cost of using the address, where 0 is free and 255 is
expensive, followed by the expiration Date after which the address should not be used,
of if null, the address never expires. After that comes a String defining the transport
protocol this router address uses. Finally there is a Mapping containing all of the
transport specific options necessary to establish the connection, such as IP address,
port number, email address, URL, etc.


+----+----+----+----+----+----+----+----+
|cost|           expiration
+----+----+----+----+----+----+----+----+
     |        transport_style           |
+----+----+----+----+-//-+----+----+----+
|                                       |
+                                       +
|               options                 |
~                                       ~
~                                       ~
|                                       |
+----+----+----+----+----+----+----+----+

cost :: Integer
        length -> 1 byte

        case 0 -> free
        case 255 -> expensive

expiration :: Date (must be all zeros, see notes below)
              length -> 8 bytes

              case null -> never expires

transport_style :: String
                   length -> 1-256 bytes

options :: Mapping
*/

// RouterAddress is the represenation of an I2P RouterAddress.
//
// https://geti2p.net/spec/common-structures#routeraddress
type RouterAddress struct {
	cost            *Integer
	expiration      *Date
	transport_style *I2PString
	options         *Mapping
	parserErr       error
}

// Bytes returns the router address as a []byte.
func (router_address RouterAddress) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, router_address.cost.Bytes()...)
	bytes = append(bytes, router_address.expiration.Bytes()...)
	strData, err := router_address.transport_style.Data()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("RouterAddress.Bytes: error getting transport_style bytes")
	} else {
		bytes = append(bytes, strData...)
	}
	//bytes = append(bytes, router_address.options.Bytes()...)
	return bytes
}

// Cost returns the cost for this RouterAddress as a Go integer.
func (router_address RouterAddress) Cost() int {
	return router_address.cost.Int()
}

// Expiration returns the expiration for this RouterAddress as an I2P Date.
func (router_address RouterAddress) Expiration() Date {
	return *router_address.expiration
}

// TransportStyle returns the transport style for this RouterAddress as an I2PString.
func (router_address RouterAddress) TransportStyle() I2PString {
	return *router_address.transport_style
}

// Options returns the options for this RouterAddress as an I2P Mapping.
func (router_address RouterAddress) Options() Mapping {
	return *router_address.options
}

//
// Check if the RouterAddress is empty or if it is too small to contain valid data.
//
func (router_address RouterAddress) checkValid() (err error, exit bool) {
	/*addr_len := len(router_address)
	exit = false
	if addr_len == 0 {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) checkValid",
			"reason": "no data",
		}).Error("invalid router address")
		err = errors.New("error parsing RouterAddress: no data")
		exit = true
	} else if addr_len < ROUTER_ADDRESS_MIN_SIZE {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) checkValid",
			"reason": "data too small (len < ROUTER_ADDRESS_MIN_SIZE)",
		}).Warn("router address format warning")
		err = errors.New("warning parsing RouterAddress: data too small")
	}*/
	if router_address.parserErr != nil {
		exit = true
	}
	return
}

// ReadRouterAddress returns RouterAddress from a []byte.
// The remaining bytes after the specified length are also returned.
// Returns a list of errors that occurred during parsing.
func ReadRouterAddress(data []byte) (router_address RouterAddress, remainder []byte, err error) {
	if data == nil || len(data) == 0 {
		log.WithField("at", "(RouterAddress) ReadRouterAddress").Error("no data")
		err = errors.New("error parsing RouterAddress: no data")
		router_address.parserErr = err
		return
	}
	cost, remainder, err := NewInteger([]byte{data[0]}, 1)
	router_address.cost = cost
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) ReadNewRouterAddress",
			"reason": "error parsing cost",
		}).Warn("error parsing RouterAddress")
		router_address.parserErr = err
	}
	expiration, remainder, err := NewDate(remainder)
	router_address.expiration = expiration
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) ReadNewRouterAddress",
			"reason": "error parsing expiration",
		}).Error("error parsing RouterAddress")
		router_address.parserErr = err
	}
	transport_style, remainder, err := NewI2PString(remainder)
	router_address.transport_style = transport_style
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) ReadNewRouterAddress",
			"reason": "error parsing transport_style",
		}).Error("error parsing RouterAddress")
		router_address.parserErr = err
	}
	options, remainder, errs := NewMapping(remainder)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) ReadNewRouterAddress",
			"reason": "error parsing options",
		}).Error("error parsing RouterAddress")
		router_address.parserErr = err
	}
	router_address.options = options
	if err != nil {
		log.WithFields(log.Fields{
			"at":     "(RouterAddress) ReadNewRouterAddress",
			"reason": "error parsing options",
		}).Error("error parsing RouterAddress")
		router_address.parserErr = err
	}
	return
}

// NewRouterAddress creates a new *RouterAddress from []byte using ReadRouterAddress.
// Returns a pointer to RouterAddress unlike ReadRouterAddress.
func NewRouterAddress(data []byte) (router_address *RouterAddress, remainder []byte, err error) {
	objrouteraddress, remainder, err := ReadRouterAddress(data)
	router_address = &objrouteraddress
	return
}
