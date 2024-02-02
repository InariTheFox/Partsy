namespace Partsy.Shared.EventBus.Abstractions
{
    public interface IRequestHandler<in TRequest, TResponse>
        where TRequest : IRequest<TResponse>
    {
        Task<TResponse> Handle(TRequest request, CancellationToken cancellationToken);
    }

    public interface RequestHandler<in TRequest>
        where TRequest : IRequest
    {
        Task Handle(TRequest request, CancellationToken cancellationToken);
    }
}