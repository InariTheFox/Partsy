namespace Partsy.Shared.EventBus.Abstractions
{
    public class Unit : IEquatable<Unit>, IComparable<Unit>, IComparable
    {
        private static readonly Unit _value = new();

        public static ref readonly Unit Value => ref _value;

        public int CompareTo(Unit other) => 0;

        public bool Equals(Unit other) => true;

        public override int GetHashCode() => 0;

        public override string ToString() => "()";

        public static bool operator ==(Unit first, Unit second) => true;

        public static bool operator !=(Unit first, Unit second) => false;

        int IComparable.CompareTo(object? obj) => 0;
    }
}