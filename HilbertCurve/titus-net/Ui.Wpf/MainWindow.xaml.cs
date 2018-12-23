using System.Windows;
using Ui.Wpf.Controls;

namespace Ui.Wpf
{
    public partial class MainWindow : Window
    {
        public MainWindow()
        {
            InitializeComponent();

            this.iteration1.Content = new TabView(new Business.Iteration1.HilbertCurve());
            this.iteration2.Content = new TabView(new Business.Iteration2.Element(7));
        }
    }
}
