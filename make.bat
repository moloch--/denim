@echo off
SETLOCAL

SET VERSION=0.0.1
SET CMDPKG=github.com/moloch--/denim/cmd
SET LDFLAGS="-s -w -X %CMDPKG%.Version=%VERSION%"

@echo on
go build -trimpath -ldflags %LDFLAGS% -o denim.exe .
@echo off

ENDLOCAL