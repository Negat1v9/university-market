package static

import (
	"fmt"

	eventmodel "github.com/Negat1v9/work-marketplace/model/event"
	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
)

var (
	ErrorPreCheckOutAmount = "Что-то пошло не так"
	NoTgCmdFinded          = "<b>Хмм...</b> что то пошло не так"
	ErrBot                 = "Кажется что-то пошло не так. Попробуйте позже"
	FilesExpected          = "Пришлите <b>файл</b>\nЕсли у вас есть <u>только фото</u>, то вы можете прислать его как файл.\n<b>Максимум 5 файлов для одной задачи</b>"
	WaitingFiles           = "Пришлите файл (<u>только одно в сообщении</u>) для работы в чат либо опубликуйте работу сейчас.\n\nФайлы будут видны <u>всем работникам</u>, которые откликнулись на вашу работу.\n<b>Максимум 5 файлов для одной задачи</b>"
	SuccessPublichTask     = "Работа успешно опубликована!\nОжидайте отклики от работников."
	HelpCmd                = "Как это работает:"
	CreateEvent            = "Опишите эвент. Если эвент должен содержать фото, то пришлите <u>фото в чат</u>, с описание эвента в <b>одном сообщении</b>.\n\nЕсли эвент не должен содержать фото, пришлите <u>только текст</u>"
	OnStartSendEvent       = "Бот начал отсылать сообщения пользователям… "
	WorkerSelected         = "🥳 Мы рады сообщить вам, что пользователь выбрал вас для выполнения работы.\nПожалуйста, свяжитесь с ним для обсуждения деталей."
	ErrTypeOnCreateEvent   = "Пришлите текст или фото с описание эвента"
	ErrNoUserName          = "Чтобы поделился контактом с пользователем у вас должен быть <u>username</u>.\nЕго можно создать в настройках профиля <b>телеграм</b>."
)

func MsgStartCommand(nickname string) string {
	return fmt.Sprintf("🥳 <b>%s</b> добро пожаловать в @shopStudentBot\n\n🌐 Здесь вы можете опубликовать работу, для института и найти работника, который ее сделает.\n\n💯 Также вы можете сами стать работником и зарабатывать на работах других студентов\n\n<b>Поддерживай свой учебный процесс на высоте и получайте доход, помогая другим</b>", nickname)
}
func MsgErrOnSuccessPayment(tgPaymentID string) string {
	return fmt.Sprintf("Кажется с нашей стороны произошла ошибка\n<b>Обратитесь пожалуйста в поддержку</b>\nПерешлите это сообщение в чат поддержки, мы обязательно разберемся с вашей проблемой.\n<code>%s</code>", tgPaymentID)
}
func BalancePayment(amount int) string {
	return fmt.Sprintf("Баланс вашего аккаунта успешно пополнен на %d звезд", amount)
}

func OnRespondFromWorker() string {
	return "🎉 <b>Новый отклик!</b>\nМожете посмотреть профиль работника если вас все устраивает <b>выберите</b> этого специалиста"
}

func AddInformationTask(meta *taskmodel.TaskMeta) string {
	return fmt.Sprintf("\n\n<b>ℹ️ О работе:</b>\n⚪️ <u>Институт</u>: %s\n⚪️ <u>Задание</u>: %s\n⚪️ <u>Предмет</u>: %s\n", meta.University, meta.TaskType, meta.Subject)
}
func SuccessAttachFiles() string {
	return fmt.Sprintf("<b>Файл успешно добавлен</b>\n%s", WaitingFiles)
}

func MsgShareContact(username string) string {
	return fmt.Sprintf("Вы выбрали <code>%s</code> работника для своей работы, можете связаться с ним для уточнения деталей! @%s\n\n❗️ <b>Если работник пытается совершить какую либо незаконную операцию, сделайте репорт на этого работника. Это можно сделать на странице с информацией о работнике.</b>", username, username)
}

func MsgEventInfoForAdmin(e *eventmodel.Event) string {
	return fmt.Sprintf("Вы создали эвент!\nИнформация:\n\nДля типа пользователя: <b>%s</b>\nОписание: <b>%s</b>\nC изображением: <b>%t</b>\n\nЕсли вы создали эвент с изображение то пришлите его в чат, если без, то нажмите кнопку начать рассылку", e.UserType, e.Caption, e.WithImage)
}
