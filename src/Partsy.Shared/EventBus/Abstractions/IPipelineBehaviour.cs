namespace Partsy.Shared.EventBus.Abstractions
{
    public delegate Task<TResponse> RequestHandlerDelegate<TResponse>();

    public interface IPipelineBehaviour<in TRequest, TResponse>
        where TRequest : notnull
    {
        Task<TResponse> Handle(TRequest request, RequestHandlerDelegate<TResponse> next, CancellationToken cancellationToken);
    }
}