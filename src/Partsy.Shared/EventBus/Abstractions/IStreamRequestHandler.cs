namespace Partsy.Shared.EventBus.Abstractions
{
    public interface IStreamRequestHandler<in TRequest, out TResponse>
        where TRequest: IStreamRequest<TResponse>
    {
        IAsyncEnumerable<TResponse> Handle(TRequest request, CancellationToken cancellationToken);
    }
}