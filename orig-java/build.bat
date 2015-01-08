@echo off
setlocal

if not exist javacpath.txt (
  echo error: please put the path to javac.exe in javacpath.txt
  goto :eof
)
set /p JAVAC=<javacpath.txt

for %%f in (lib\*.jar) do (
  rem echo %%f
  call set cp=%%cp%%;..\%%f
)

for /r src %%f in (*.class) do (
  del "%%f"
)

rem echo %cp%
cd src
%JAVAC% -cp .%cp% org\stathissideris\ascii2image\core\CommandLineConverter.java 2> tmpx
cd ..

more src\tmpx

endlocal
