using System;
using System.Collections.Generic;
using System.Linq;

namespace Business.Iteration2
{
    public class Element : HilbertCurve
    {
        #region Attributes

        #endregion

        #region Custom Properties

        #endregion

        #region Properties

        public Dictionary<Position, Element> Elements { get; }

        public Orientation Orientation { get; set; }

        private Int32 Childrens { get; set; }

        #endregion

        #region Constructor

        public Element(Int32 childrens = 2, Orientation orientation = Orientation.Down)
        {
            this.Elements = new Dictionary<Position, Element>(4)
            {
                { Position.Zero, default(Element) },
                { Position.One, default(Element) },
                { Position.Two, default(Element) },
                { Position.Three, default(Element) }
            };

            this.Childrens = childrens;

            this.Orientation = orientation;

            this.Populate(childrens);
        }

        #endregion

        #region Methods

        public List<Coordinates> GetCoordenades(Coordinates basePosition = default(Coordinates), Int32 level = 0)
        {
            this.UpdateChildrenOrientation();

            if (this.Elements[Position.Zero] == null)
                return new List<Coordinates>
                {
                    this.PositionZero() + basePosition,
                    this.PositionOne() + basePosition,
                    this.PositionTwo() + basePosition,
                    this.PositionThree() + basePosition
                };

            var list = new List<Coordinates>();

            var separator = (int) Math.Pow(2, Childrens);
            
            list.AddRange(this.Elements[Position.Zero].GetCoordenades(this.PositionZero() * separator + basePosition, level + 1));

            list.AddRange(this.Elements[Position.One].GetCoordenades(this.PositionOne() * separator + basePosition, level + 1));

            list.AddRange(this.Elements[Position.Two].GetCoordenades(this.PositionTwo() * separator + basePosition, level + 1));

            list.AddRange(this.Elements[Position.Three].GetCoordenades(this.PositionThree() * separator + basePosition, level + 1));

            return list;
        }

        private void Populate(Int32 childrens)
        {
            if (childrens <= 0) return;

            this.Elements[Position.Zero] = new Element(childrens - 1);
            this.Elements[Position.One] = new Element(childrens - 1);
            this.Elements[Position.Two] = new Element(childrens - 1);
            this.Elements[Position.Three] = new Element(childrens - 1);
        }

        private Coordinates PositionZero()
        {
            switch (this.Orientation)
            {
                case Orientation.Up:
                case Orientation.Left:
                    return new Coordinates(1, 1);

                case Orientation.Down:
                case Orientation.Right:
                    return new Coordinates(0, 0);

                default:
                    throw new NotImplementedException();
            }
        }

        private Coordinates PositionOne()
        {
            switch (this.Orientation)
            {
                case Orientation.Up:
                case Orientation.Right:
                    return new Coordinates(1, 0);

                case Orientation.Down:
                case Orientation.Left:
                    return new Coordinates(0, 1);

                default:
                    throw new NotImplementedException();
            }
        }

        private Coordinates PositionTwo()
        {
            switch (this.Orientation)
            {
                case Orientation.Up:
                case Orientation.Left:
                    return new Coordinates(0, 0);

                case Orientation.Down:
                case Orientation.Right:
                    return new Coordinates(1, 1);

                default:
                    throw new NotImplementedException();
            }
        }

        private Coordinates PositionThree()
        {
            switch (this.Orientation)
            {
                case Orientation.Up:
                case Orientation.Right:
                    return new Coordinates(0, 1);

                case Orientation.Down:
                case Orientation.Left:
                    return new Coordinates(1, 0);

                default:
                    throw new NotImplementedException();
            }
        }

        private void UpdateChildrenOrientation()
        {
            if (this.Elements[Position.Zero] != null)
            {
                if (this.Orientation == Orientation.Up) this.Elements[Position.Zero].Orientation = Orientation.Left;
                if (this.Orientation == Orientation.Left) this.Elements[Position.Zero].Orientation = Orientation.Up;
                if (this.Orientation == Orientation.Down) this.Elements[Position.Zero].Orientation = Orientation.Right;
                if (this.Orientation == Orientation.Right) this.Elements[Position.Zero].Orientation = Orientation.Down;
            }

            if (this.Elements[Position.One] != null) this.Elements[Position.One].Orientation = this.Orientation;

            if (this.Elements[Position.Two] != null) this.Elements[Position.Two].Orientation = this.Orientation;

            if (this.Elements[Position.Three] != null)
            {
                if (this.Orientation == Orientation.Up) this.Elements[Position.Three].Orientation = Orientation.Right;
                if (this.Orientation == Orientation.Left) this.Elements[Position.Three].Orientation = Orientation.Down;
                if (this.Orientation == Orientation.Down) this.Elements[Position.Three].Orientation = Orientation.Left;
                if (this.Orientation == Orientation.Right) this.Elements[Position.Three].Orientation = Orientation.Up;
            }
        }

        public override IEnumerable<Coordinates> GetCoordenades()
        {
            return this.GetCoordenades();
        }

        #endregion
    }
}
