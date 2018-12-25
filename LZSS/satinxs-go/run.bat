@ECHO OFF

ECHO Bigger compression ration == better > results.log
ECHO ------------------------ >> results.log
ECHO . >> results.log
go run baseline.go >> results.log

ECHO ------------------------ >> results.log
ECHO . >> results.log
CD iteration2
ECHO Iteration 2 >> ..\results.log
go run main.go bitio.go  >> ..\results.log
CD ..

ECHO ------------------------ >> results.log
ECHO . >> results.log
CD iteration3
ECHO Iteration 3 >> ..\results.log
go build main.go bitio.go encoder.go decoder.go symbol.go 
ECHO Best configuration turned to be 22 16 4. To run configuration tester, run without arguments. >> ..\results.log
main 22 16 4 >> ..\results.log  
rm main.exe
CD ..

ECHO . >> results.log
ECHO . >> results.log
ECHO DONE >> results.log