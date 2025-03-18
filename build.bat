@echo off
REM Script to build local-share for Windows

REM Create bin directory if it doesn't exist
if not exist bin mkdir bin

REM Build the single binary
echo Building local-share...
go build -o bin\local-share.exe .\cmd

if %ERRORLEVEL% == 0 (
  echo Build successful! Binary located at bin\local-share.exe
  echo.
  echo Usage examples:
  echo   Start server:       .\bin\local-share.exe receiver
  echo   Send text message:  .\bin\local-share.exe send text ^<server-ip^> "message"
  echo   Send file:          .\bin\local-share.exe send file ^<server-ip^> C:\path\to\file
  echo   Show help:          .\bin\local-share.exe help
) else (
  echo Build failed!
) 