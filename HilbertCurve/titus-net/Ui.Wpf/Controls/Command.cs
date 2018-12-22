using System;
using System.Windows.Input;

namespace Ui.Wpf.Controls
{
    public class Command<T> : ICommand
        where T : class
    {
        private readonly Predicate<T> canExecute;
        private readonly Action<T> execute;

        public Command(Action<T> execute)
            : this(execute, null)
        {
        }

        public Command(Action<T> execute, Predicate<T> canExecute)
        {
            this.execute = execute;
            this.canExecute = canExecute;
        }

        public Boolean CanExecute(Object parameter)
        {
            if (canExecute == null)
                return true;

            return canExecute((T)parameter);
        }

        public void Execute(Object parameter)
        {
            execute((T)parameter);
        }

        public event EventHandler CanExecuteChanged;
        public void RaiseCanExecuteChanged()
        {
            if (CanExecuteChanged != null)
                CanExecuteChanged(this, EventArgs.Empty);
        }
    }
}
