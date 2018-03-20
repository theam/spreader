package main

import (
	"encoding/json"
	"fmt"
	"github.com/alvaroloes/kinesis/rtc"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/siddontang/go-mysql/canal"
	"log"
)


type MySQLEvent struct {
	Database string
	Table string
	Type string
	OldValues map[string]interface{}
	NewValues map[string]interface{}
}

func newMySQLEvent(e *canal.RowsEvent, oldValues, newValues map[string]interface{}) MySQLEvent {
	return MySQLEvent{
		Database:  e.Table.Schema,
		Table:     e.Table.Name,
		Type:      e.Action,
		OldValues: oldValues,
		NewValues: newValues,
	}
}

type RowEventProcessor struct {
	canal.DummyEventHandler
	bus *kinesis.Kinesis
}

func (processor *RowEventProcessor) OnRow(e *canal.RowsEvent) error {
	fmt.Println("New SQL event. Type: ", e.Action)
	// Format the changed rows as events
	events := formatMySQLEvents(e)
	// Encode them as JSON
	encodedEvents, _ := json.MarshalIndent(events, "", "\t")
	// Send them to Kinesis
	res, err := processor.bus.PutRecord(&kinesis.PutRecordInput{
		StreamName:   aws.String(rtc.StreamName),
		Data:         encodedEvents,
		PartitionKey: aws.String("partition-key"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print the result for confirmation
	fmt.Println(res)
	return nil
}


func formatMySQLEvents(e *canal.RowsEvent) []MySQLEvent {
	var events []MySQLEvent
	switch e.Action {
	case canal.InsertAction:
		for _, rowValues := range e.Rows {
			event := newMySQLEvent(e, nil, formatRow(e, rowValues))
			events = append(events, event)
		}
	case canal.UpdateAction:
		// When updating rows, MySQL sends columns in pairs: {old value, new value}
		for i := 0; i < len(e.Rows); i += 2 {
			oldRowValues := e.Rows[i]
			newRowValues := e.Rows[i+1]
			event := newMySQLEvent(e, formatRow(e, oldRowValues), formatRow(e, newRowValues))
			events = append(events, event)
		}
	case canal.DeleteAction:
		for _, rowValues := range e.Rows {
			event := newMySQLEvent(e, formatRow(e, rowValues), nil)
			events = append(events, event)
		}
	}
	return events
}

func formatRow(e *canal.RowsEvent, values []interface{}) map[string]interface{} {
	formattedRow := make(map[string]interface{})
	for i, column := range e.Table.Columns {
		formattedRow[column.Name] = values[i]
	}
	return formattedRow
}
