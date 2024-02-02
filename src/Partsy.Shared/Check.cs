namespace Partsy.Shared
{
    public static class Check
    {
        public static T IsTrue<T>(bool condition, Func<T> lambda, Exception failureException)
        {
            if (!condition)
            {
                throw failureException;
            }

            return lambda();
        }

        public static T NotNull<T>(T value, string parameterName)
        {
            if (null == value)
            {
                throw new ArgumentNullException(parameterName);
            }

            return value;
        }
    }
}