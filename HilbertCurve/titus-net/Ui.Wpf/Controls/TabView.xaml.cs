using Business;
using System;
using System.Text.RegularExpressions;
using System.Windows.Controls;

namespace Ui.Wpf.Controls
{
    public partial class TabView  : UserControl
    {
        private static readonly Regex _regex = new Regex("[^0-9.-]+"); 

        public TabView(HilbertCurve iteration) 
        {
            InitializeComponent();

            this.DataContext = new TabViewModel(iteration);
        }

        private new void PreviewTextInput(Object sender, System.Windows.Input.TextCompositionEventArgs e)
        {
            e.Handled = !IsTextAllowed(e.Text);
        }

        private static Boolean IsTextAllowed(String text)
        {
            return !_regex.IsMatch(text);
        }
    }
}
