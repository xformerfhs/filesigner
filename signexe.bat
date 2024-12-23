@echo off
if A%1==A goto :errNoArg

filesigner.exe sign %1 --include-file filesigner --include-file filesigner.exe  --name exe-%1

call :writeLog Signing had return code %errorlevel%

goto :end

:errNoArg
call :writeLog Context id is missing

:writeLog
echo %date% %time% %*
exit /b

:end
