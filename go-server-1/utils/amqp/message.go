package amqp

/*
 * Message is the abstraction of queue message which can have two parts:
 *    	Headers : meta data of the message
 *      Body	: actual message which is produced by producer
 *		it can be acknowledge, negatively acknowledge, reject, and requeue.
 */
type Message interface {
	/*
	 * GetID retrieve message ID which is provided by different provider
	 */
	GetID() string

	/*
	 * Ack when consumer has finished work on a delivery, consumer need to acknowledge the message.
	 *      Options are the parameters which can be different for each provider (eg. rabbitmq and etc...)
	 *      		read more on specific provider implemetation.
	 */
	Ack(flag bool) error

	/*
	 * Headers Return the message headers
	 */
	Headers() map[string]interface{}

	/*
	 * Body Return the message body
	 */
	Body() []byte
}
