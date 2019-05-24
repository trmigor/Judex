install:
	mkdir -pv /var/log/Judex
	touch /var/log/Judex/logs.log
	mkdir -pv /var/www/judex.vdi.mipt.ru/html/template
	mkdir -pv /var/www/judex.vdi.mipt.ru/html/static
	mkdir -pv /var/www/judex.vdi.mipt.ru/emails
	mkdir -pv /var/www/judex.vdi.mipt.ru/problems
	mkdir -pv /var/www/judex.vdi.mipt.ru/solutions

configure:
	cp templates/* /var/www/judex.vdi.mipt.ru/html/template
	cp static/* /var/www/judex.vdi.mipt.ru/html/static
	cp emails/* /var/www/judex.vdi.mipt.ru/emails
