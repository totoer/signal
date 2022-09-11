package main

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"net"
// 	"signal/db"
// 	"signal/sip"
// 	"signal/transport"

// 	"github.com/google/uuid"
// )

// const (
// 	CLUSTER_STATE_KEY = "/cluster/state"
// 	NODE_MAP_KEY      = "/cluster/node/map/%s"
// )

// type Node struct {
// 	ID       uuid.UUID `json:"id"`
// 	Host     string    `json:"host"`
// 	Port     int       `josn:"port"`
// 	Capacity int       `josn:"capacity"`
// }

// func (n *Node) Send(t transport.Transport, m sip.Message) error {
// 	// m.GetHeaders().Append("Route", sip.URI{})
// 	addr := net.UDPAddr{
// 		IP:   net.ParseIP(n.Host),
// 		Port: n.Port,
// 	}
// 	return t.Send(&addr, m.Data())
// }

// type ClusterState struct {
// 	Nodes    []Node `json:"nodes"`
// 	Pointer  int    `json:"pointer"`
// 	Capacity int    `json:"capacity"`
// }

// type Cluster struct {
// 	db *db.DB
// }

// func (c *Cluster) AppendNode(ctx context.Context, n Node) error {
// 	var err error = nil

// 	err = c.db.Lock(ctx, CLUSTER_STATE_KEY, func() {
// 		var cs *ClusterState
// 		if err1 := c.db.Get(ctx, CLUSTER_STATE_KEY, cs); err1 != nil {
// 			err = err1
// 		} else {
// 			cs.Nodes = append(cs.Nodes, n)
// 			err = c.db.Put(ctx, CLUSTER_STATE_KEY, cs)
// 		}
// 	})

// 	return err
// }

// var ErrEmptyCluster = errors.New("empty cluster")

// func (c *Cluster) GetNextNode(ctx context.Context) (*Node, error) {
// 	var node *Node = nil
// 	var err error = nil

// 	err = c.db.Lock(ctx, CLUSTER_STATE_KEY, func() {
// 		var cs *ClusterState
// 		if err := c.db.Get(ctx, CLUSTER_STATE_KEY, cs); err != nil {
// 			if len(cs.Nodes) == 0 {
// 				err = ErrEmptyCluster
// 			} else if len(cs.Nodes) == 1 {
// 				node = &cs.Nodes[0]
// 			} else {
// 				node = &cs.Nodes[cs.Pointer]
// 				if cs.Capacity > 0 {
// 					cs.Capacity -= 1
// 				} else {
// 					cs.Pointer = (cs.Pointer + 1) % len(cs.Nodes)
// 					cs.Capacity = cs.Nodes[cs.Pointer].Capacity
// 				}
// 				c.db.Put(ctx, CLUSTER_STATE_KEY, cs)
// 			}
// 		}
// 	})

// 	return node, err
// }

// var ErrKeyNotPresent = errors.New("key not present")
// var ErrNodeNotPresent = errors.New("node not present")
// var ErrWrongNodeID = errors.New("wrong node id")

// func (c *Cluster) GetNode(ctx context.Context, cid sip.PlainHeader, m *sip.MethodType) (*Node, error) {
// 	key := fmt.Sprintf(NODE_MAP_KEY, string(cid.Value))
// 	if nodeID, err := c.db.GetString(ctx, key); err != nil {
// 		if m != nil && m.IncludeIn(sip.INVITE, sip.REGISTER, sip.OPTIONS) {
// 			if node, err := c.GetNextNode(ctx); err != nil {
// 				return nil, err
// 			} else {
// 				return node, nil
// 			}
// 		} else {
// 			return nil, ErrNodeNotPresent
// 		}
// 	} else {
// 		var cs *ClusterState
// 		if err := c.db.Get(ctx, CLUSTER_STATE_KEY, cs); err != nil {
// 			return nil, err
// 		} else {
// 			for _, node := range cs.Nodes {
// 				if node.ID.String() == nodeID {
// 					return &node, nil
// 				}
// 			}
// 		}
// 		return nil, ErrWrongNodeID
// 	}
// }
