namespace Partsy.Shared.EventBus.Abstractions
{
    public interface IRequest : IBaseRequest { }

    public interface IRequest<out TResponse> : IBaseRequest { }

    public interface IBaseRequest
    {
        Guid Id { get; }

        DateTimeOffset CreatedAt { get; }
    }
}