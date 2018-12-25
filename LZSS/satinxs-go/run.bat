@echo off

echo Bigger compression ration == better > results.log
echo ------------------------ >> results.log
echo . >> results.log
go run baseline.go >> results.log

echo ------------------------ >> results.log
echo . >> results.log
cd iteration2
echo Iteration 2 >> ..\results.log
go run main.go bitio.go  >> ..\results.log
cd ..

echo ------------------------ >> results.log
echo . >> results.log
cd iteration3
echo Iteration 3 >> ..\results.log
go run main.go bitio.go encoder.go decoder.go symbol.go >> ..\results.log
cd ..

echo . >> results.log
echo . >> results.log
echo DONE >> results.log