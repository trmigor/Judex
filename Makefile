install:
	mkdir -pv /var/log/Judex
	touch /var/log/Judex/logs.log
	mkdir -pv /var/www/judex.vdi.mipt.ru/html/template
	mkdir -pv /var/www/judex.vdi.mipt.ru/html/static
	mkdir -pv /var/www/judex.vdi.mipt.ru/emails
	mkdir -pv /var/www/judex.vdi.mipt.ru/problems
	mkdir -pv /var/www/judex.vdi.mipt.ru/solutions

configure:
	cp -r templates/* /var/www/judex.vdi.mipt.ru/html/template
	cp -r static/* /var/www/judex.vdi.mipt.ru/html/static
	cp -r emails/* /var/www/judex.vdi.mipt.ru/emails
