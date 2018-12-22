using Business;
using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Input;
using System.Windows.Media;

namespace Ui.Wpf.Controls
{
    public class TabViewModel : ViewModel
    {
        #region Attributes

        private readonly HilbertCurve iteration;

        private Int32 miliseconds;

        #endregion

        #region Properties

        private IList<Coordinates> CoordintesList { get; set; }

        #endregion

        #region Blinded Properties
        
        public PointCollection Points
        {
            get
            {
                return new PointCollection(CoordintesList.Select(coordinates => new Point(10 + coordinates.X * 10, 1 + coordinates.Y * 10)));
            }
        }

        public String CoordinatesLog
        {
            get
            {
                return String.Join(Environment.NewLine, this.CoordintesList);
            }
        }

        public Int32 Miliseconds
        {
            get
            {
                return this.miliseconds;
            }
            set
            {
                this.miliseconds = value;

                this.RaisePropertyChanged("Miliseconds");
            }
        }

        private ICommand drawCommand;
        public ICommand DrawCommand
        {
            get
            {
                return this.drawCommand;
            }
        }

        private ICommand clearCommand;
        public ICommand ClearCommand
        {
            get
            {
                return this.clearCommand;
            }
        }
        
        #endregion

        #region Constructor

        public TabViewModel(HilbertCurve iteration)
        {
            this.iteration = iteration;

            this.CoordintesList = new List<Coordinates>();
            
            this.RaisePropertyChanged("Points");
        }

        #endregion

        #region Methods

        public void AddPoint(Coordinates coordinates)
        {
            this.CoordintesList.Add(coordinates);

            this.RaisePropertyChanged("Points");
            this.RaisePropertyChanged("CoordinatesLog");
        }

        public void DeleteAllPoints()
        {
            this.CoordintesList.Clear();

            this.RaisePropertyChanged("Points");
            this.RaisePropertyChanged("CoordinatesLog");
        }

        protected override void AddCommands()
        {
            base.AddCommands();

            this.drawCommand = new Command<String>(s => this.DrawAsync());

            this.clearCommand = new Command<String>(s => this.Clear());
        }

        private async void DrawAsync()
        {
            this.Clear();

            var progress = new Progress<Coordinates>(this.AddPoint);

            await Task.Run(() => this.Draw(progress));
        }

        private void Draw(IProgress<Coordinates> progress)
        {
            foreach (var coordinates in this.iteration.DoSomething())
            {
                if (this.Miliseconds > 0) Thread.Sleep(this.Miliseconds);

                this.AddPoint(coordinates);
            }
        }

        private void Clear()
        {
            this.DeleteAllPoints();
        }

        #endregion
    }
}
