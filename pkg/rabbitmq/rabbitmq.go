package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	channel  *amqp.Channel    // 与 RabbitMQ 的通道连接
	conn     *amqp.Connection // 与 RabbitMQ 的连接
	Name     string           // 队列名称
	exchange string           // 绑定的交换机名称
}

func New(s string) *RabbitMQ {
	connection, err := amqp.Dial(s)
	if err != nil {
		panic(err)
	}
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	// 声明 匿名队列
	queue, err := channel.QueueDeclare(
		"",
		false, // 非持久化队列
		true,  // 队列在没有消费者时是否自动删除， true 表示会自动删除
		false, // 是否是排他队列， 排他队列，只能由当前队列连接使用，其他连接不能访问
		false, // 是否等待队列的声明结果
		nil,
	)

	if err != nil {
		panic(err)
	}

	mq := new(RabbitMQ)
	mq.channel = channel
	mq.conn = connection
	mq.Name = queue.Name
	return mq
}

// 绑定交换机
func (q *RabbitMQ) Bind(exchange string) {
	err := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	q.exchange = exchange
}

// 发送消息到指定的队列
func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	// 发送给指定的队列
	err = q.channel.Publish("", queue, false, false, amqp.Publishing{
		ReplyTo: q.Name,
		Body:    []byte(str),
	})
	if err != nil {
		panic(err)
	}
}

// 发送消息到指定的 exchange 上
func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	err = q.channel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		},
	)
	if err != nil {
		panic(err)
	}

}

// 消费队列的消息
func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	consum, err := q.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	return consum
}

func (q *RabbitMQ) Close() {
	q.channel.Close()
	q.conn.Close()
}
