using System;

namespace Business
{
    public struct Coordinates
    {
        #region Properties

        public Int32 X { get; set; }

        public Int32 Y { get; set; }

        #endregion

        #region Constructor

        public Coordinates(Int32 x, Int32 y) : this()
        {
            this.X = x;

            this.Y = y;
        }

        #endregion

        #region Operators

        public static Coordinates operator *(Coordinates coordinates, Int32 number)
        {
            return new Coordinates(coordinates.X * number, coordinates.Y * number);
        }

        public static Coordinates operator +(Coordinates left, Coordinates right)
        {
            return new Coordinates(left.X + right.X, left.Y + right.Y);
        }

        #endregion

        #region Methods

        public override String ToString()
        {
            return String.Format("( {0} , {1} )", this.X, this.Y);
        }

        #endregion
    }
}
