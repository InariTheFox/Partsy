namespace Partsy.Shared.EventBus.RabbitMQ
{
    public class RabbitMQConnectionSettings
    {
        public string HostName { get; set; } = "localhost";

        public string Password { get; set; } = "guest";

        public int Port { get; set; } = 5672;

        public int RetryCount { get; set; } = 5;

        public string UserName { get; set; } = "guest";

        public string VirtualHost { get; set; } = "/";
    }
}