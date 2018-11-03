/* Copyright (c) 2018 Gabor Seljan
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

// Package mikrotik implements methods to interact with RouterOS API.
//
// Available functionality is currently limited to adding or removing
// an IP address to or from an address list. API credentials should be
// stored in the EveBox Server configuration file.
package mikrotik

import (
	"github.com/jasonish/evebox/log"
	"github.com/spf13/viper"
	"gopkg.in/routeros.v2"
	"errors"
)

// Connect to RouterOS API.
func Dial() (*routeros.Client, error) {
	address := viper.GetString("mikrotik.address")
	username := viper.GetString("mikrotik.username")
	password := viper.GetString("mikrotik.password")

	if viper.GetBool("mikrotik.tls") == true {
		return routeros.DialTLS(address, username, password, nil)
	}

	return routeros.Dial(address, username, password)
}

// Send sentence to RouterOS API.
func Exec(command ...string) (*routeros.Reply, error) {
	c, err := Dial()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if viper.GetBool("mikrotik.async") == true {
		c.Async()
	}

	log.Debug("Executing %s", command)
	r, err := c.RunArgs(command)
	if err != nil {
		log.Debug("%v", err)
	}

	return r, err
}

// Get internal ID of IP address for removal.
func GetIPAddressID(address string) (string, error) {
	list := viper.GetString("mikrotik.list")
	id := ""

	r, err := Exec(
		"/ip/firewall/address-list/print",
		"?list="+list,
		"?address="+address,
		"=.proplist=.id")

	if err != nil {
		return id, errors.New("Failed to retrive ID of IP address")
	}

	for _, re := range r.Re {
		id = re.Map[".id"]
	}

	return id, nil
}

// Add the specified IP address to the pre-configured address list.
func AddIPAddressToList(address string, comment string) error {
	list := viper.GetString("mikrotik.list")

	_, err := Exec(
		"/ip/firewall/address-list/add",
		"=list="+list,
		"=address="+address,
		"=comment="+comment,
		"=disabled=no")

	if err != nil {
		err = errors.New("Failed to add IP address to list")
	}

	return err
}

// Remove the specified IP address from the pre-configured address list.
func RemoveIPAddressFromList(address string) error {
	id, err := GetIPAddressID(address)

	if err != nil && id == "" {
		return err
	}

	_, err = Exec(
		"/ip/firewall/address-list/remove",
		"=.id="+id)

	if err != nil {
		err = errors.New("Failed to remove IP address from list")
	}

	return err
}
