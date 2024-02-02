using RabbitMQ.Client;

namespace Partsy.Shared.EventBus.RabbitMQ
{
    public interface IPersistentRabbitMQConnection : IDisposable
    {
        bool IsConnected { get; }

        IModel CreateModel();

        bool TryConnect();
    }
}