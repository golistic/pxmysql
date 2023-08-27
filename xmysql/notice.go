// Copyright (c) 2022, 2023, Geert JM Vanderkelen

package xmysql

import (
	"fmt"

	"github.com/golistic/pxmysql/internal/mysqlx/mysqlxnotice"
	"github.com/golistic/pxmysql/xmysql/internal/network"
)

type StateChanges struct {
	ClientID          uint64
	GeneratedInsertID uint64
	RowsAffected      uint64
	CurrentSchema     string
	ProducedMessage   string
}

type notices struct {
	warnings                     []*mysqlxnotice.Warning
	sessionVariableChanges       []*mysqlxnotice.SessionVariableChanged
	sessionStateChanges          []*mysqlxnotice.SessionStateChanged
	groupReplicationStateChanges []*mysqlxnotice.GroupReplicationStateChanged
	serverHello                  *mysqlxnotice.ServerHello
	unhandled                    []mysqlxnotice.Frame_Type
	stateChanges                 StateChanges
}

func (n *notices) add(msg *network.ServerMessage) error {
	frame := &mysqlxnotice.Frame{}
	if err := msg.Unmarshall(frame); err != nil {
		return fmt.Errorf("failed unmarshalling notice message (%w)", err)
	}

	switch mysqlxnotice.Frame_Type(frame.GetType()) {
	case mysqlxnotice.Frame_WARNING:
		m := &mysqlxnotice.Warning{}
		if err := msg.Unmarshall(m); err != nil {
			return err
		}
		n.warnings = append(n.warnings, m)
	case mysqlxnotice.Frame_SESSION_VARIABLE_CHANGED:
		m := &mysqlxnotice.SessionVariableChanged{}
		if err := msg.Unmarshall(m); err != nil {
			return fmt.Errorf("failed unmarshalling '%s' (%w)", m.String(), err)
		}
		n.sessionVariableChanges = append(n.sessionVariableChanges, m)
	case mysqlxnotice.Frame_SESSION_STATE_CHANGED:
		m := &mysqlxnotice.SessionStateChanged{}
		if err := network.UnmarshalPartial(frame.Payload, m); err != nil {
			return fmt.Errorf("failed unmarshalling '%s' (%w)", m.String(), err)
		}
		network.Trace("state", m)

		switch m.GetParam() {
		case mysqlxnotice.SessionStateChanged_GENERATED_INSERT_ID:
			if len(m.Value) > 0 {
				n.stateChanges.GeneratedInsertID = m.Value[0].GetVUnsignedInt()
			}
		case mysqlxnotice.SessionStateChanged_ROWS_AFFECTED:
			if len(m.Value) > 0 {
				n.stateChanges.RowsAffected = m.Value[0].GetVUnsignedInt()
			}
		case mysqlxnotice.SessionStateChanged_CURRENT_SCHEMA:
			if len(m.Value) > 0 {
				n.stateChanges.CurrentSchema = string(m.Value[0].VString.Value)
			}
		case mysqlxnotice.SessionStateChanged_PRODUCED_MESSAGE:
			if len(m.Value) > 0 {
				n.stateChanges.ProducedMessage = string(m.Value[0].VString.Value)
			}
		case mysqlxnotice.SessionStateChanged_CLIENT_ID_ASSIGNED:
			if len(m.Value) > 0 {
				n.stateChanges.ClientID = m.Value[0].GetVUnsignedInt()
			}
		}

		n.sessionStateChanges = append(n.sessionStateChanges, m)
	case mysqlxnotice.Frame_GROUP_REPLICATION_STATE_CHANGED:
		m := &mysqlxnotice.GroupReplicationStateChanged{}
		if err := msg.Unmarshall(m); err != nil {
			return fmt.Errorf("failed unmarshalling '%s' (%w)", m.String(), err)
		}
		n.groupReplicationStateChanges = append(n.groupReplicationStateChanges, m)
	case mysqlxnotice.Frame_SERVER_HELLO:
		m := &mysqlxnotice.ServerHello{}
		if err := msg.Unmarshall(m); err != nil {
			return fmt.Errorf("failed unmarshalling '%s' (%w)", m.String(), err)
		}
		n.serverHello = m
	default:
		n.unhandled = append(n.unhandled, mysqlxnotice.Frame_Type(frame.GetType()))
	}

	return nil
}
