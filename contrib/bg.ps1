#! pwsh

Start-Process  `
    -FilePath C:\WINDOWS\system32\cmd.exe  `
    -ArgumentList "start `"`" /b cmd /c .\chaindb.x64-mswindows.exe 1>>chaindb.log 2>&1"  `
    -WindowStyle Hidden
