using System.Net.Sockets;
using Microsoft.Extensions.Logging;
using Polly;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using RabbitMQ.Client.Exceptions;

namespace Partsy.Shared.EventBus.RabbitMQ
{
    internal class DefaultPersistentRabbitMQConnection : IPersistentRabbitMQConnection
    {
        private readonly IConnectionFactory _connectionFactory;
        private readonly ILogger _logger;
        private readonly int _retryCount;

        private IConnection? _connection;
        private bool _disposed;
        private object _sync = new();

        public DefaultPersistentRabbitMQConnection(RabbitMQConnectionSettings settings, ILogger<DefaultPersistentRabbitMQConnection> logger)
        {
            Check.NotNull(settings, nameof(settings));

            _connectionFactory = new ConnectionFactory()
            {
                HostName = settings.HostName,
                Port = settings.Port,
                UserName = settings.UserName,
                Password = settings.Password,
                VirtualHost = settings.VirtualHost
            };

            _retryCount = settings.RetryCount;

            _logger = logger;
        }

        public bool IsConnected
            => null != _connection && _connection.IsOpen && !_disposed;

        public IModel CreateModel()
            => Check.IsTrue(
                IsConnected,
                () => _connection!.CreateModel(),
                new InvalidOperationException("No RabbitMQ connections are available to perform this action."));

        public void Dispose()
        {
            if (_disposed)
            {
                return;
            }

            _disposed = true;

            try
            {
                _connection?.Dispose();
            }
            catch (IOException ex)
            {
                _logger.LogCritical(ex, "Failed to dispose of RabbitMQ connection.");
            }
        }

        public bool TryConnect()
        {
            _logger.LogInformation("Attempting to connect to RabbitMQ");

            lock (_sync)
            {
                var policy = Policy.Handle<SocketException>()
                    .Or<BrokerUnreachableException>()
                    .WaitAndRetry(
                        _retryCount,
                        retryAttempt => TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)),
                        (ex, time) =>
                        {
                            _logger.LogWarning(ex, "Could not connect to RabbitMQ after {TimeOut}s ({ExceptionMessage})", $"{time.TotalSeconds:n1}", ex.Message);
                        });

                policy.Execute(() =>
                {
                    _connection = _connectionFactory.CreateConnection();
                });

                if (IsConnected)
                {
                    _connection!.CallbackException += OnCallbackException;
                    _connection!.ConnectionBlocked += OnConnectionBlocked;
                    _connection!.ConnectionShutdown += OnConnectionShutdown;

                    return true;
                }
                else
                {
                    _logger.LogCritical("Fatal exception: RabbitMQ connection could not be created.");

                    return false;
                }
            }
        }

        private void OnCallbackException(object? sender, CallbackExceptionEventArgs args)
        {
            if (_disposed)
            {
                return;
            }

            _logger.LogWarning(args.Exception, "A RabbitMQ connection has thrown an exception. See inner exception for details. Attempting to reconnect.");

            TryConnect();
        }

        private void OnConnectionBlocked(object? sender, ConnectionBlockedEventArgs args)
        {
            if (_disposed)
            {
                return;
            }

            _logger.LogWarning("A RabbitMQ connection has been blocked. Attempting to reconnect.");

            TryConnect();
        }

        private void OnConnectionShutdown(object? sender, ShutdownEventArgs args)
        {
            if (_disposed)
            {
                return;
            }

            _logger.LogWarning(args.Exception, "A RabbitMQ connection has been shutdown. See inner exception for details if any. Attempting to reconnect.");

            TryConnect();
        }
    }
}