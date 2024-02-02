using System.Reflection.Metadata;

namespace Partsy.Shared.EventBus.Abstractions
{
    public interface INotificationHandler<in TNotification>
        where TNotification : INotification
    {
        Task Handle(TNotification notification, CancellationToken cancellationToken);
    }

    public abstract class NotificationHandler<TNotification> : INotificationHandler<TNotification>
        where TNotification : INotification
    {
        Task INotificationHandler<TNotification>.Handle(TNotification notification, CancellationToken cancellationToken)
        {
            Handle(notification);

            return Task.CompletedTask;
        }

        protected abstract void Handle(TNotification notification);
    }
}