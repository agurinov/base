import os
import sys
import json

import pygeoip


def get_geo_record(ip):
	try:
		# get geo database and
		# get record by ip address ( dict )
		geo_record = pygeoip.GeoIP(
			os.path.join(
				os.getcwd(),
				'geoip_city.dat'
			)
		).record_by_addr(ip)
	except (pygeoip.GeoIPError, IOError, AttributeError):
		geo_record = None

	geo_record = geo_record or {}
	geo_record['ip'] = ip

	return geo_record


if __name__ == '__main__':
	ip = sys.stdin.read()
	print json.dumps(
		get_geo_record(ip)
	)
