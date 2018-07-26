from django.shortcuts import render, render_to_response
from django.http import HttpResponse, Http404, HttpResponseRedirect
from AdminPanel.models import AdminPanelUser, AdminPanelAd
import hashlib, simplejson
from django.conf import settings
from datetime import datetime, timezone, timedelta

def now():
    return datetime.now(timezone.utc) + timedelta(minutes=180)


def index(request):
    if 'admin' not in request.session:
        return render(request, 'AdminPanel/Authorization/wrapper.html')
    else:
        result = {}
        id = request.session['admin']
        obj_db = AdminPanelUser.objects.filter(id=id)[0]
        result['first_name'], result['last_name'] = obj_db.first_name, obj_db.last_name
        result['mail'] = obj_db.mail
        result['username'] = obj_db.username
        result['status'] = obj_db.status

        print(obj_db.id)

        if obj_db.status == '7':
            result['allow_admin_table'] = True

        ads = AdminPanelAd.objects.all()

        result['ads'] = []

        for ad in ads:
            if ad.date_start < now() < ad.date_stop:
                admin_id_tmp = ad.admin_id
                admin_tmp = AdminPanelUser.objects.filter(id=admin_id_tmp)[0]
                result['ads'].append({'id': ad.id, 'text': ad.ad_text, 'admin_first_name': admin_tmp.first_name,
                                  'admin_second_name': admin_tmp.last_name, 'admin_mail': admin_tmp.mail})

        return render(request, 'AdminPanel/MainPage/wrapper.html', result)

def authorization(request):
    if 'admin' not in request.session:
        if request.POST:
            user_name = request.POST['login']
            password = request.POST['password']

            sha256 = hashlib.sha256()
            md5 = hashlib.md5()

            sha256.update(password.encode('utf-8'))
            md5.update(sha256.hexdigest().encode('utf-8'))
            enc_pass = md5.hexdigest()

            if AdminPanelUser.objects.filter(username=user_name):
                admin_usr = AdminPanelUser.objects.filter(username=user_name)[0]
                if admin_usr.password == enc_pass:
                    request.session['admin'] = admin_usr.id
                    request.session.modified = True
                    data = {'authorization': True, 'type': 0} # успешная авторизация
                else:
                    data = {'authorization': False, 'type': 1} # неправильный пароль
            else:
                data = {'authorization': False, 'type': 2} # не существует такого пользователя
        else:
            raise Http404
    else:
        data = {'authorization': False, 'type': 3} # пользователь уже авторизован!

    result = simplejson.dumps(data)
    print(data)
    return HttpResponse(result)

def out(request):
    if 'admin' in request.session:
        del request.session['admin']
        return HttpResponseRedirect(settings.HOSTNAME + "custom_admin/")
    return HttpResponseRedirect(settings.HOSTNAME + "404/")