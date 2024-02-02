using System.Net.Sockets;
using System.Text;
using System.Text.Json;
using Microsoft.Extensions.Logging;
using Partsy.Shared.EventBus.Abstractions;
using Polly;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using RabbitMQ.Client.Exceptions;

namespace Partsy.Shared.EventBus.RabbitMQ
{
    public class RabbitMQBus : IEventBus
    {
        const string brokerName = "partsy_event_bus";
        const string queueName = "default_queue";

        private readonly Dictionary<string, ICollection<Type>> _requestHandlerRegistrations = new();

        private readonly AutoResetEvent _autoReset = new(false);
        private readonly IPersistentRabbitMQConnection _connection;
        private readonly ILogger _logger;
        private readonly int _retryCount;

        private IModel _consumerChannel;

        public RabbitMQBus(
            IPersistentRabbitMQConnection connection,
            RabbitMQConnectionSettings settings,
            ILogger<RabbitMQBus> logger)
        {
            _connection = Check.NotNull(connection, nameof(connection));
            _logger = Check.NotNull(logger, nameof(logger));
            _retryCount = settings.RetryCount;
        }

        public IAsyncEnumerable<TResponse> CreateStream<TResponse>(IStreamRequest<TResponse> request, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }

        public Task Publish<TNotification>(TNotification message, CancellationToken cancellationToken = default)
            where TNotification : INotification
        {
            throw new NotImplementedException();
        }

        public Task<TResponse> Publish<TResponse>(IRequest<TResponse> request, CancellationToken cancellationToken = default)
        {
            if (!_connection.IsConnected)
            {
                _connection.TryConnect();
            }

            var policy = Policy.Handle<SocketException>()
                .Or<BrokerUnreachableException>()
                .WaitAndRetry(
                    _retryCount,
                    retryAttempt => TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)),
                    (ex, time) =>
                    {
                        _logger.LogWarning(ex, "Could not publish message to RabbitMQ: {EventId} after {Timeout}s. See inner exception for details", request.Id, $"{time.TotalSeconds}:n1");
                    });

            var responseName = typeof(TResponse).Name;

            InternalSubscribe(queueName, responseName);

            var requestName = request.GetType().Name;

            _logger.LogTrace("Creating RabbitMQ channel to send request: {EventId} ({RequestName})", request.Id, requestName);

            TResponse? response = default;

            using (var channel = _connection.CreateModel())
            {
                _logger.LogTrace("Declaring RabbitMQ exchange to send request: {EventId}", request.Id);

                channel.ExchangeDeclare(brokerName, "direct", true, true, null);

                var consumer = new EventingBasicConsumer(channel);
                consumer.Received += (ch, ea) =>
                {
                    var responseBody = Encoding.UTF8.GetString(ea.Body.ToArray());

                    try
                    {
                        response = JsonSerializer.Deserialize<TResponse?>(responseBody);

                        channel.BasicAck(ea.DeliveryTag, false);
                    }
                    catch (Exception ex)
                    {
                        _logger.LogError(ex, "Unable to deserialize response from RabbitMQ. See inner exception for details.");
                    }
                };

                var message = JsonSerializer.Serialize(request);
                var body = Encoding.UTF8.GetBytes(message);

                policy.Execute(() =>
                {
                    var props = channel.CreateBasicProperties();
                    props.DeliveryMode = 2;

                    _logger.LogTrace("Sending request to RabbitMQ: {EventId}", request.Id);

                    channel.BasicPublish(
                        exchange: brokerName,
                        routingKey: requestName,
                        mandatory: true,
                        basicProperties: props,
                        body: body
                    );

                    channel.BasicConsume(queueName, false, consumer);
                });
            }

            return Task.FromResult(response);
        }

        public Task Send<TRequest>(TRequest request, CancellationToken cancellationToken = default)
            where TRequest : IRequest
        {
            if (!_connection.IsConnected)
            {
                _connection.TryConnect();
            }

            var policy = Policy.Handle<SocketException>()
                .Or<BrokerUnreachableException>()
                .WaitAndRetry(
                    _retryCount,
                    retryAttempt => TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)),
                    (ex, time) =>
                    {
                        _logger.LogWarning(ex, "Could not send message to RabbitMQ: {EventId} after {Timeout}s. See inner exception for details.", request.Id, $"{time.TotalSeconds}:n1");
                    });

            var requestName = request.GetType().Name;

            _logger.LogTrace("Creating RabbitMQ channel to send request: {EventId} ({RequestName})", request.Id, requestName);

            using (var channel = _connection.CreateModel())
            {
                _logger.LogTrace("Declaring RabbitMQ exchange to send request: {EventId}", request.Id);

                channel.ExchangeDeclare(brokerName, "direct", true, true, null);

                var message = JsonSerializer.Serialize(request);
                var body = Encoding.UTF8.GetBytes(message);

                policy.Execute(() =>
                {
                    var props = channel.CreateBasicProperties();
                    props.DeliveryMode = 2;

                    _logger.LogTrace("Sending request to RabbitMQ: {EventId}", request.Id);

                    channel.BasicPublish(
                        exchange: brokerName,
                        routingKey: requestName,
                        mandatory: true,
                        basicProperties: props,
                        body: body
                    );
                });
            }

            return Task.CompletedTask;
        }

        public Task<TResponse> Send<TResponse>(IRequest<TResponse> request, CancellationToken cancellationToken = default)
        {
            if (!_connection.IsConnected)
            {
                _connection.TryConnect();
            }

            var policy = Policy.Handle<SocketException>()
                .Or<BrokerUnreachableException>()
                .WaitAndRetry(
                    _retryCount,
                    retryAttempt => TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)),
                    (ex, time) =>
                    {
                        _logger.LogWarning(ex, "Could not publish message to RabbitMQ: {EventId} after {Timeout}s. See inner exception for details", request.Id, $"{time.TotalSeconds}:n1");
                    });

            var responseName = typeof(TResponse).Name;

            InternalSubscribe(queueName, responseName);

            var requestName = request.GetType().Name;

            _logger.LogTrace("Creating RabbitMQ channel to send request: {EventId} ({RequestName})", request.Id, requestName);

            TResponse? response = default;

            using (var channel = _connection.CreateModel())
            {
                _logger.LogTrace("Declaring RabbitMQ exchange to send request: {EventId}", request.Id);

                channel.ExchangeDeclare(brokerName, "direct", true, true, null);

                var consumer = new EventingBasicConsumer(channel);
                consumer.Received += (ch, ea) =>
                {
                    var responseBody = Encoding.UTF8.GetString(ea.Body.ToArray());

                    try
                    {
                        response = JsonSerializer.Deserialize<TResponse?>(responseBody);

                        channel.BasicAck(ea.DeliveryTag, false);
                    }
                    catch (Exception ex)
                    {
                        _logger.LogError(ex, "Unable to deserialize response from RabbitMQ. See inner exception for details.");
                    }
                };

                var message = JsonSerializer.Serialize(request);
                var body = Encoding.UTF8.GetBytes(message);

                policy.Execute(() =>
                {
                    var props = channel.CreateBasicProperties();
                    props.DeliveryMode = 2;

                    _logger.LogTrace("Sending request to RabbitMQ: {EventId}", request.Id);

                    channel.BasicPublish(
                        exchange: brokerName,
                        routingKey: requestName,
                        mandatory: true,
                        basicProperties: props,
                        body: body
                    );

                    channel.BasicConsume(queueName, false, consumer);
                });
            }

            return Task.FromResult(response);
        }

        public Task Subscribe<TRequest, TRequestHandler, TResponse>(string queueName, CancellationToken cancellationToken = default)
             where TRequest : IRequest<TResponse>
             where TRequestHandler : IRequestHandler<TRequest, TResponse>
        {
            var requestName = typeof(TRequest).Name;

            InternalSubscribe(queueName, requestName);

            if (_requestHandlerRegistrations.TryGetValue(requestName, out var handlers))
            {
                handlers.Add(typeof(TRequestHandler));
            }
            else
            {
                _requestHandlerRegistrations.Add(requestName, [typeof(TRequestHandler)]);
            }

            return Task.CompletedTask;
        }

        public Task Unsubscribe<TRequest, TRequestHandler, TResponse>(CancellationToken cancellationToken = default)
            where TRequest : IRequest<TResponse>
            where TRequestHandler : IRequestHandler<TRequest, TResponse>
        {
            var requestName = typeof(TRequest).Name;

            if (_requestHandlerRegistrations.TryGetValue(requestName, out var handlers))
            {
                if (handlers.Contains(typeof(TRequestHandler)))
                {
                    handlers.Remove(typeof(TRequestHandler));
                }

                if (handlers.Count == 0)
                {
                    _requestHandlerRegistrations.Remove(requestName);
                }
            }

            return Task.CompletedTask;
        }

        private void InternalSubscribe(string queueName, string routingKey)
        {
            if (!_connection.IsConnected)
            {
                _connection.TryConnect();
            }

            _consumerChannel = _connection.CreateModel();
            _logger.LogTrace("Creating Queue {QueueName}", queueName);

            _consumerChannel.QueueDeclare(queueName, true, false, false, null);

            _logger.LogTrace("Binding Queue {QueueName} to {ExchangeName} using {RoutingKey}", queueName, brokerName, routingKey);

            _consumerChannel.QueueBind(queueName, brokerName, routingKey, null);
        }

        private void StartConsume(string queueName)
        {
            _logger.LogTrace("Starting RabbitMQ consumer on channel {QueueName}", queueName);

            var consumer = new AsyncEventingBasicConsumer(_consumerChannel);

            consumer.Received += OnConsumerReceived;

            _consumerChannel.BasicConsume(queueName, false, consumer);
        }

        private async Task OnConsumerReceived(object sender, BasicDeliverEventArgs args)
        {
            var requestName = args.RoutingKey;
            var message = Encoding.UTF8.GetString(args.Body.ToArray());

            try
            {
                // TODO: Process message.
            }
            catch (Exception ex)
            {
                _logger.LogWarning(ex, "Error processing message. See inner exception for details.");
            }

            ((IModel)sender).BasicAck(args.DeliveryTag, multiple: false);

            _autoReset.Set();
        }

        private void WaitOne()
        {
            while (!_autoReset.WaitOne())
            {}
        }
    }
}