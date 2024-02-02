namespace Parsty.Shared.EventBus.Abstractions
{
    public delegate IAsyncEnumerable<TResponse> StreamHandlerDelegate<TResponse>();

    public interface IStreamPipelineBehaviour<in TRequest, TResponse>
        where TResponse : notnull
    {
        IAsyncEnumerable<TResponse> Handle(TRequest request, StreamHandlerDelegate<TResponse> next, CancellationToken cancellationToken);
    }
}