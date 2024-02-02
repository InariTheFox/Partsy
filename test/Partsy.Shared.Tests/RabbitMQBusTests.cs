using Microsoft.Extensions.Logging;
using Partsy.Shared.EventBus.Abstractions;
using Partsy.Shared.EventBus.RabbitMQ;

namespace Partsy.Shared.Tests
{
    [TestClass]
    public class RabbitMQBusTests
    {
        private readonly ILoggerFactory _loggerFactory;

        public RabbitMQBusTests()
        {
            _loggerFactory = LoggerFactory.Create((config) =>
            {
                config.AddConsole().SetMinimumLevel(LogLevel.Trace);
            });
        }

        [TestMethod]
        public void ConstructorWithNullParametersTest()
        {
            Assert.ThrowsException<ArgumentNullException>(() => new RabbitMQBus(null, null, null));
        }

        [TestMethod]
        public void ConstructorTest()
        {
            var settings = new RabbitMQConnectionSettings();
            var logger = _loggerFactory.CreateLogger<RabbitMQBus>();
            var connectionFactory = new PersistentRabbitMQConnectionFactory(settings, _loggerFactory);
            var connection = connectionFactory.CreateConnection();

            Assert.IsNotNull(connection);

            var bus = new RabbitMQBus(connection, settings, logger);

            Assert.IsNotNull(bus);
        }

        [TestMethod]
        public void SendTest()
        {
            var settings = new RabbitMQConnectionSettings
            {
                HostName = "10.10.0.40",
                UserName = "test",
                Password = "test",
                VirtualHost = "test"
            };

            var logger = _loggerFactory.CreateLogger<RabbitMQBus>();
            var connectionFactory = new PersistentRabbitMQConnectionFactory(settings, _loggerFactory);
            var connection = connectionFactory.CreateConnection();

            Assert.IsNotNull(connection);

            var bus = new RabbitMQBus(connection, settings, logger);

            Assert.IsNotNull(bus);

            Exception? exception = null;

            try
            {
                bus.Send(new SimpleRequest());
            }
            catch (Exception ex)
            {
                exception = ex;
            }

            Assert.IsNull(exception);
        }

        [TestMethod]
        public void SubscribeTest()
        {
            var settings = new RabbitMQConnectionSettings
            {
                HostName = "10.10.0.40",
                UserName = "test",
                Password = "test",
                VirtualHost = "test"
            };

            var logger = _loggerFactory.CreateLogger<RabbitMQBus>();
            var connectionFactory = new PersistentRabbitMQConnectionFactory(settings, _loggerFactory);
            var connection = connectionFactory.CreateConnection();

            Assert.IsNotNull(connection);

            var bus = new RabbitMQBus(connection, settings, logger);

            Assert.IsNotNull(bus);

            Exception? exception = null;

            try
            {
                var response = bus.Send(new SimpleRequestWithResponse());

                Assert.IsNotNull(response, "Response is null.");
                Assert.IsInstanceOfType<Task<SimpleResponse>>(response, $"Response is not of Type {typeof(Task<SimpleResponse>)}.");
                Assert.IsTrue(response.IsCompleted, "Response is not completed.");
                Assert.IsNotNull(response.Result, "Result is null.");
            }
            catch (Exception ex)
            {
                exception = ex;
            }

            Assert.IsNull(exception, exception?.Message);
        }

        private class SimpleRequest : IRequest
        {
            public SimpleRequest()
            {
                Id = Guid.NewGuid();
                CreatedAt = DateTimeOffset.Now;
            }

            public Guid Id { get; }

            public DateTimeOffset CreatedAt { get; }

        }

        private class SimpleRequestWithResponse : IRequest<SimpleResponse>
        {
            public SimpleRequestWithResponse()
            {
                Id = Guid.NewGuid();
                CreatedAt = DateTimeOffset.Now;
            }

            public Guid Id { get; }

            public DateTimeOffset CreatedAt { get; }

            public bool Handled { get; set; }
        }

        private class SimpleResponse
        {
            public bool Handled { get; set; } = false;
        }

        private class SimpleRequestHandler : RequestHandler<SimpleRequest>
        {
            public Task Handle(SimpleRequest request, CancellationToken cancellationToken)
            {
                return Task.FromResult(Unit.Value);
            }
        }

        private class SimpleRequestWithResponseHandler : IRequestHandler<SimpleRequestWithResponse, SimpleResponse>
        {
            public Task<SimpleResponse> Handle(SimpleRequestWithResponse request, CancellationToken cancellationToken)
            {
                return Task.FromResult(new SimpleResponse { Handled = true });
            }
        }
    }
}