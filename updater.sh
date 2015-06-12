#!/bin/bash
/etc/init.d/ottemo stop
#backup old foundation by date
cp /home/ottemo/foundation /home/ottemo/foundation/backup/foundation-$(date +%Y%m%d)
#copy latest foundation just scp'd into place
cp /home/ottemo/deploy/foundation-latest /home/ottemo/foundation
/etc/init.d/ottemo start
