@echo off 
SETLOCAL 
pushd %%~dp0 
for /D %%%%D in (services\*) do @echo %%%%D 
popd 
ENDLOCAL 
