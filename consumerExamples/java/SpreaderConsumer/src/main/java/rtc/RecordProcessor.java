package rtc;

import com.amazonaws.services.dynamodbv2.model.SourceTableDetails;
import com.amazonaws.services.kinesis.AmazonKinesis;
import com.amazonaws.services.kinesis.model.*;

import java.nio.charset.Charset;
import java.util.List;

public class RecordProcessor implements Runnable {
    private AmazonKinesis client;
    private String streamName;
    private final String shardId;

    public RecordProcessor(AmazonKinesis client, String streamName, String shardId) {
        this.client = client;
        this.streamName = streamName;
        this.shardId = shardId;
    }
    public void run() {
        GetShardIteratorRequest getShardIteratorRequest = new GetShardIteratorRequest();
        getShardIteratorRequest.setStreamName(this.streamName);
        getShardIteratorRequest.setShardId(shardId);
        getShardIteratorRequest.setShardIteratorType("LATEST");

        GetShardIteratorResult getShardIteratorResult = client.getShardIterator(getShardIteratorRequest);
        String shardIterator = getShardIteratorResult.getShardIterator();

        while (true) {
            // Create a new getRecordsRequest with an existing shardIterator
            // Set the maximum records to return to 25
            GetRecordsRequest getRecordsRequest = new GetRecordsRequest();
            getRecordsRequest.setShardIterator(shardIterator);
            getRecordsRequest.setLimit(25);

            GetRecordsResult result = client.getRecords(getRecordsRequest);

            printRecords(result.getRecords());

            try {
                Thread.sleep(1000);
            }
            catch (InterruptedException exception) {
                throw new RuntimeException(exception);
            }

            shardIterator = result.getNextShardIterator();
        }
    }

    private void printRecords(List<Record> records) {
        if (records.size() == 0) {
            return;
        }

        System.out.println("The following records were found:");
        for (Record record : records) {
            byte[] bytes = record.getData().array();
            String data = new String( bytes, Charset.forName("UTF-8") );
            System.out.println("Published at: " + record.getApproximateArrivalTimestamp() + "Data: " + data);
        }
    }
}
