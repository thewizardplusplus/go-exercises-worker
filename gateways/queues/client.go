package queues

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Client ...
type Client struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewClient ...
func NewClient(queueDSN string, queueMaximalSize int) (Client, error) {
	connection, err := amqp.Dial(queueDSN)
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to dial the message broker")
	}

	channel, err := connection.Channel()
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to open the channel")
	}

	if err := channel.Qos(queueMaximalSize, 0, false); err != nil {
		return Client{}, errors.Wrap(err, "unable to set the queue maximal size")
	}

	for _, queueName := range []string{
		SolutionQueueName,
		SolutionResultQueueName,
	} {
		if _, err := channel.QueueDeclare(
			queueName, // queue name
			true,      // durable
			false,     // auto-delete
			false,     // exclusive
			false,     // no wait
			nil,       // arguments
		); err != nil {
			return Client{}, errors.Wrapf(err, "unable to declare queue %q", queueName)
		}
	}

	client := Client{connection: connection, channel: channel}
	return client, nil
}

// Close ...
func (client Client) Close() error {
	if err := client.channel.Close(); err != nil {
		return errors.Wrap(err, "unable to close the channel")
	}

	if err := client.connection.Close(); err != nil {
		return errors.Wrap(err, "unable to close the connection")
	}

	return nil
}
