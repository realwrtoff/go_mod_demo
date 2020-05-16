#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import getopt
from urllib import parse
import requests


def main():
    opts, args = getopt.getopt(sys.argv, "h", ["help"])
    if 2 > len(args):
        action = 'xx'
    else:
        action = args[1]
    url = 'http://127.0.0.1:7060/{0}'.format(action)
    if action == 'channel':
        advertiser_url = 'http://127.0.0.1:21688/click'
        params = {
            'pub': 'didazhuan',
            'cid': 'ddz_xxx',
            'status': 1,
            'url': advertiser_url,
            'origin_cid': 'guahao',
            'name': 'haina'
        }
    elif action == 'click':
        callback_url = 'http://127.0.0.1:21688/install?f=中国'
        params = {
            'pub': 'didazhuan',
            'cid': 'ddz_xxx',
            'ip': '192.168.1.111',
            'devicetype': 'iphone',
            'os': 'ios',
            'osversion': 'IOS-1.13',
            'idfa': 'fengmin-de-ipnone',
            'callback': callback_url
        }
    elif action == 'install':
        params = {
            'click_id': '5ebfdfeae5994eba616dab91'
        }
    else:
        print('action {0} not support'.format(action))
        sys.exit(-1)

    print('request {0}'.format(url))
    print('params {0}'.format(params))
    r = requests.get(url, params=params)
    print('response status{0}, content [{1}]'.format(r.status_code, r.text))


if __name__ == '__main__':
    main()
