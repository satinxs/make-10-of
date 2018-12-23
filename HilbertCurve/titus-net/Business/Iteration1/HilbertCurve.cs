using System;
using System.Collections.Generic;

namespace Business.Iteration1
{
    public class HilbertCurve : Business.HilbertCurve
    {
        public override IEnumerable<Coordinates> GetCoordenades()
        {
            return new List<Coordinates>()
            {
                new Coordinates(0, 0),
                new Coordinates(0, 1),
                new Coordinates(1, 1),
                new Coordinates(1, 0),
                new Coordinates(2, 0),
                new Coordinates(3, 0),
                new Coordinates(3, 1),
                new Coordinates(2, 1),
                new Coordinates(2, 2),
                new Coordinates(3, 2),
                new Coordinates(3, 3),
                new Coordinates(2, 3),
                new Coordinates(1, 3),
                new Coordinates(1, 2),
                new Coordinates(0, 2),
                new Coordinates(0, 3)
            };
        }
    }
}
