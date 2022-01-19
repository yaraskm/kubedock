//import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import org.testcontainers.containers.KafkaContainer;
import org.testcontainers.containers.Network;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.utility.DockerImageName;

@Testcontainers
//@Disabled
public class KafkaTest {
    
    @Test
    void testKafkaStartup() {
        Network network = Network.newNetwork();

        // https://docs.confluent.io/5.5.0/release-notes/index.html supports Apache Kafka version 2.5.0
        KafkaContainer KAFKA_CONTAINER = new KafkaContainer(DockerImageName.parse("confluentinc/cp-kafka:5.5.0"))
                .withStartupAttempts(3)
                .withNetworkAliases("foo")
                .withNetwork(network);
        
        KAFKA_CONTAINER.start();
        KAFKA_CONTAINER.stop();
    }
}
