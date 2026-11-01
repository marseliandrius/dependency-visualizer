@echo off
setlocal
pushd "%~dp0"
if errorlevel 1 exit /b 1
go run ./src %*
set "runExitCode=%errorlevel%"
popd
exit /b %runExitCode%
