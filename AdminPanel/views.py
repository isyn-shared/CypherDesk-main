from django.shortcuts import render
from django.http import HttpResponse, Http404, HttpResponseRedirect
from AdminPanel.models import AdminPanelUser
import hashlib, simplejson
from django.conf import settings

def index(request):
    if 'admin' not in request.session:
        return render(request, 'AdminPanel/Authorization/wrapper.html')
    else:
        return render(request, 'AdminPanel/MainPage/wrapper.html')

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