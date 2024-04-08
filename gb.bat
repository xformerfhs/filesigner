@echo off
rem Build for go programs
rem If an argument is provided the executables are compressed
setlocal

call :writeLog Start %0

if exist D:\MyPrograms\upx\upx.exe set upxPath=D:\MyPrograms\upx\upx.exe
if exist "C:\Program Files\UPX\upx.exe" set upxPath=C:\Program Files\UPX\upx.exe
set useCompression=N
if "%1" == "" goto :startProc
if not exist "%upxPath%" goto :noUPX
set useCompression=Y
goto :startProc
:noUPX
call :writeLog UPX does not exists at %upxPath%

:startProc
set CGO_ENABLED=0
set GOARCH=amd64
set GOAMD64=v2

call :getExecutableName

echo.
set GOOS=windows
call :buildExecutable
echo.

set GOOS=linux
call :buildExecutable
echo.

goto :endProc

:getExecutableName
set thisPath=%~p0
if "%thisPath:~-1%" == "\" set thisPath=%thisPath:~0,-1%
for %%f in ("%thisPath%") do set executableName=%%~nxf
exit /b

:buildExecutable
call :writeLog Start build for %GOOS%/%GOARCH%
set realExecutableName=%executableName%
if %GOOS% == windows set realExecutableName=%realExecutableName%.exe
rem -s: Omit the symbol table and debug information
rem -w: Omit the DWARF symbol table
rem -trimpath: Remove all file system paths from the resulting executable
go build -ldflags="-s -w" -trimpath
set rc=%errorlevel%
call :writeLog Go build for %GOOS%/%GOARCH% has error level %rc%
if %rc% GTR 0 exit /b
if %useCompression% == Y "%upxPath%" --best --lzma -v %realExecutableName%
call :writeLog Built %GOOS%/%GOARCH% executable %realExecutableName%
exit /b

:writeLog
echo %date% %time% %*
exit /b

:endProc
call :writeLog End %0
