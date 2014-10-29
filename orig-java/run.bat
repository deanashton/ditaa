@echo off
setlocal

if not exist javapath.txt (
  echo error: please put the path to java.exe in javapath.txt
  goto :eof
)
set /p JAVA=<javapath.txt

for %%f in (lib\*.jar) do (
  call set cp=%%cp%%;%%f
)

%JAVA% -cp src%cp% org.stathissideris.ascii2image.core.CommandLineConverter %*

endlocal
