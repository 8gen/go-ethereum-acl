// Copyright 2014 The go-ethereum Authors
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
	"bufio"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	datadirACL = "acl"
)


type Role int64

const (
    Sender Role = iota
    Recipient
    Creator
)

func (role Role) GoString() string {
    switch role {
    case Sender:
        return "Sender"
    case Recipient:
        return "Recipient"
    default:
        return "Creator"
    }
}

// Who is allowed to transfer
type ACL struct {
    config *Config
    done chan bool
    senders []string `toml:",omitempty"`
    recipients []string `toml:",omitempty"`
    creators []string `toml:",omitempty"`
}

// Is creator permitted?
func (acl *ACL) permitted(role Role) []string {
    switch role {
    case Sender:
        return acl.senders
    case Recipient:
        return acl.recipients
    default:
        return acl.creators
    }
}

func NewACL(config *Config) *ACL {
    acl := ACL{
        config: config,
    }
    go acl.StartWatcher()
    return &acl
}

func IsEqual(a1 []string, a2 []string) bool {
   sort.Strings(a1)
   sort.Strings(a2)
   if len(a1) == len(a2){
      for i, v := range a1 {
         if (v != a2[i]){
            return false
         }
      }
   } else {
      return false
   }
   return true
}

func (acl *ACL) readFile(name string, current_addrs *[]string) {
    prefix, err := acl.config.ACLDirConfig() 
    addrs := []string{}
    file, err := os.Open(fmt.Sprintf("%s/%s", prefix, name))
    defer file.Close()
	if err != nil {
        *current_addrs = []string{}
		return 
	}

    scanner := bufio.NewScanner(file)
    // optionally, resize scanner's capacity for lines over 64K, see next example
    for scanner.Scan() {
        addr := scanner.Text()
        if len(addr) != 42 {
            log.Debug(fmt.Sprintf("During import from %s found wrong address: %s (%d)", name, addr, len(addr)))
        } else {
            addrs = append(addrs, addr)
        }
    }

    sort.Strings(addrs)

    if !reflect.DeepEqual(addrs, *current_addrs) {
        log.Debug(fmt.Sprintf("Refresh addresses from: %s, values: %#v", name, addrs))
        *current_addrs = addrs
    }
}


func (acl *ACL) readAll() (bool, error) {
    acl.readFile("allowed_from.txt", &acl.senders)
    acl.readFile("allowed_to.txt", &acl.recipients)
    acl.readFile("allowed_smart_deploy.txt", &acl.creators)
    return true, nil
}

func (acl *ACL) StartWatcher() {
    ticker := time.NewTicker(1 * time.Second)
    acl.done = make(chan bool)
    go func() {
        for {
            select {
            case <-acl.done:
                return
            case <-ticker.C:
                acl.readAll()
            }
        }
    }()
}

func (acl *ACL) StopWatcher() {
    acl.done <- true
}

func (acl *ACL) isPermitted(role Role, address common.Address) bool {
    addrs := acl.permitted(role)
    addr := strings.ToLower(address.Hex())
    flag := false;
    if len(addrs) == 0 {
        flag = true
    }
    for i := 0; i < len(addrs); i++ {
        if addrs[i] == addr {
            flag = true
        }
    }
    log.Debug(fmt.Sprintf("Check is %s permitted for %s - %v", addr, role.GoString(), flag))
    return flag
}

// Is sender permitted?
func (acl *ACL) SenderPermitted(address common.Address) bool {
    return acl.isPermitted(Sender, address)
}

// Is recipients permitted?
func (acl *ACL) RecipientPermitted(address common.Address) bool {
    return acl.isPermitted(Recipient, address)
}

// Is creator permitted?
func (acl *ACL) CreatorPermitted(address common.Address) bool {
    return acl.isPermitted(Creator, address)
}
