@echo off
:: https://www.gamedev.net/forums/topic/475776-how-do-you-combine-static-library-files/
:: You'll need to invoke vsvars32.bat first so that your enviroment variables are set correctly.
lib /out:libdn_deps_debug.lib ../../dependencies/lib/Debug/*.lib
lib /out:libdn_deps_release.lib ../../dependencies/lib/Release/*.lib
pause 0