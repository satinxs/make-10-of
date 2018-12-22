using System;
using System.ComponentModel;
using System.Runtime.CompilerServices;

namespace Ui.Wpf.Controls
{
    public class ViewModel : INotifyPropertyChanged
    {
        public event PropertyChangedEventHandler PropertyChanged;

        public ViewModel()
        {
            this.AddCommands();
        }

        protected virtual void AddCommands()
        {

        }

        protected void RaisePropertyChanged([CallerMemberName] String propertyName = "")
        {
            if (PropertyChanged != null)
            {
                PropertyChanged(this, new PropertyChangedEventArgs(propertyName));
            }
        }
    }
}
