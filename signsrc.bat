@echo off
if A%1==A goto :errNoArg

filesigner.exe sign %1 --recurse --include-file *.go --include-file go* --include-file gb* --include-file *.md --include-file filesigner_sbom.json --exclude-dir .* --name source

call :writeLog Signing had return code %errorlevel%

goto :end

:errNoArg
call :writeLog Context id is missing

:writeLog
echo %date% %time% %*
exit /b

:end
