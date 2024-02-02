using Microsoft.Extensions.Logging;

namespace Partsy.Shared.EventBus.RabbitMQ
{
    public class PersistentRabbitMQConnectionFactory
    {
        private readonly ILoggerFactory _loggerFactory;
        private readonly RabbitMQConnectionSettings _settings;

        public PersistentRabbitMQConnectionFactory(RabbitMQConnectionSettings settings, ILoggerFactory loggerFactory)
        {
            _settings = Check.NotNull(settings, nameof(settings));
            _loggerFactory = Check.NotNull(loggerFactory, nameof(loggerFactory));
        }

        public IPersistentRabbitMQConnection CreateConnection()
        {
            var logger = _loggerFactory.CreateLogger<DefaultPersistentRabbitMQConnection>();

            return new DefaultPersistentRabbitMQConnection(_settings, logger);
        }
    }
}