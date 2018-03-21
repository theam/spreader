package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"fmt"
	"github.com/theam/spreader/rtc"
	"time"
)

func main() {
	bus := rtc.NewKinesis()

	// First we get the shards from the stream
	shardRes, err := bus.ListShards(&kinesis.ListShardsInput{
		StreamName: aws.String(rtc.StreamName),
	})
	if err != nil {
		log.Fatal("Failed to get stream shards", err)
	}
	fmt.Println("Found", len(shardRes.Shards), "shards. Starting processors")

	// Now, we need to read from each shard in its own goroutine
	for _, shard := range shardRes.Shards {
		go getShardRecords(bus, shard.ShardId)
	}

	// Block main thread forever. To exit hit Ctrl+C
	<-make(chan bool)
}

func getShardRecords(bus *kinesis.Kinesis, shardId *string) {
	// Get the first shard iterator. We will start reading from the latest published record
	shardIteratorRes, err := bus.GetShardIterator(&kinesis.GetShardIteratorInput{
		StreamName: aws.String(rtc.StreamName),
		ShardId: shardId,
		ShardIteratorType: aws.String(kinesis.ShardIteratorTypeLatest),

	})
	if err != nil {
		log.Fatal(err)
	}

	iterator := shardIteratorRes.ShardIterator
	for {
		recordsRes, err := bus.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: iterator,
		})
		if err != nil {
			log.Println("Failed to get records from shard", err)
		}
		printRecordsData(recordsRes.Records)

		iterator = recordsRes.NextShardIterator
		time.Sleep(time.Second)
	}
}
func printRecordsData(records []*kinesis.Record) {
	if len(records) == 0 {
		return
	}
	fmt.Println("The following records were found:")
	for _, record := range records {
		fmt.Println("Published at: ", record.ApproximateArrivalTimestamp, "Data: ", string(record.Data))
	}
}
