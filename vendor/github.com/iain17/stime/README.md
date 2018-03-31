# Server time
This simple package will try and return the precise UTC time even though the system time might not be accurate

# How?
Using Network Time Protocol [(NTP)](https://en.wikipedia.org/wiki/Network_Time_Protocol) at start your go app will download the difference in time.