# avscan-api

avscan-api is a REST API that exposes a status page and let's you upload a single file in order to check for viruses 
with ClamAV. ClamAV can run as a clamd server process in another container or the files can be scanned locally for 
testing, which takes a whole lot longer, as the clamscan command used updates it's patterns every time it is invoked.

