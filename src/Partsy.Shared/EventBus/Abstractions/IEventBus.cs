namespace Partsy.Shared.EventBus.Abstractions
{
    public interface IEventBus
    {
        IAsyncEnumerable<TResponse> CreateStream<TResponse>(IStreamRequest<TResponse> request, CancellationToken cancellationToken = default);

        Task Publish<TNotification>(TNotification message, CancellationToken cancellationToken = default)
            where TNotification : INotification;

        Task<TResponse> Publish<TResponse>(IRequest<TResponse> message, CancellationToken cancellationToken = default);

        Task Send<TRequest>(TRequest request, CancellationToken cancellationToken = default)
            where TRequest : IRequest;

        Task Subscribe<TRequest, TRequestHandler, TResponse>(string queueName, CancellationToken cancellationToken = default)
            where TRequest : IRequest<TResponse>
            where TRequestHandler : IRequestHandler<TRequest, TResponse>;

        Task Unsubscribe<TRequest, TRequestHandler, TResponse>(CancellationToken cancellationToken = default)
            where TRequest : IRequest<TResponse>
            where TRequestHandler : IRequestHandler<TRequest, TResponse>;
    }
}