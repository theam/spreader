package rtc;
import com.amazonaws.auth.DefaultAWSCredentialsProviderChain;
import com.amazonaws.services.kinesis.AmazonKinesis;
import com.amazonaws.services.kinesis.AmazonKinesisClientBuilder;
import com.amazonaws.services.kinesis.model.DescribeStreamRequest;
import com.amazonaws.services.kinesis.model.DescribeStreamResult;
import com.amazonaws.services.kinesis.model.Shard;

import java.util.List;

public class Consumer {
    public static final String STREAM_NAME = "pizza-tube";

    public static void main(String[] args) {

        AmazonKinesisClientBuilder clientBuilder = AmazonKinesisClientBuilder.standard();
        clientBuilder.setRegion("eu-west-1");
        clientBuilder.setCredentials(new DefaultAWSCredentialsProviderChain());

        AmazonKinesis client = clientBuilder.build();

        DescribeStreamRequest describeStreamRequest = new DescribeStreamRequest();
        describeStreamRequest.setStreamName(STREAM_NAME);
        DescribeStreamResult describeStreamResult = client.describeStream(describeStreamRequest);
        List<Shard> shards = describeStreamResult.getStreamDescription().getShards();
        System.out.println("Found " + shards.size() + " shards. Starting processors");

        describeStreamResult.getStreamDescription().getShards().forEach(shard -> {
            new RecordProcessor(client, STREAM_NAME, shard.getShardId()).run();
        });
    }
}