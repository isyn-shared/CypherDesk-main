from django.core.mail import BadHeaderError, EmailMultiAlternatives
from django.http import HttpResponse, HttpResponseRedirect
from django.views.decorators.csrf import csrf_exempt
from django.core.exceptions import PermissionDenied
from django.conf import settings

@csrf_exempt
def send(request):
    client_ip = str(request.META['REMOTE_ADDR'])
    if client_ip != '127.0.0.1':
        return PermissionDenied

    if request.POST:
        mail_res = True
        from_mail = request.POST['from']
        subject = request.POST['subject']
        html_content = request.POST['html_content']
        to = request.POST['to']

        if type(to) is str:
            to = [to, ]

        text_content = request.POST['text_content']

        msg = EmailMultiAlternatives(subject, text_content, from_mail, to)
        connection = msg.get_connection()
        msg.attach_alternative(html_content, "text/html")

        if subject and html_content and from_mail and to:
            try:
                msg.send()
                # message: 'Thanks, your email was sent'
            except BadHeaderError:
                # message: 'Invalid header found.'
                mail_res = False
            except AttributeError:
                # Неправильно указана почта
                mail_res = False
            except:
                mail_res = False
        else:
            # In reality we'd use a form class
            # to get proper validation errors.
            # message: Make sure all fields are entered and valid
            mail_res = False

        return HttpResponse(mail_res)
    return HttpResponseRedirect(settings.HOSTNAME + '404/')