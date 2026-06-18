import telebot
from telebot.types import InlineKeyboardMarkup, InlineKeyboardButton, WebAppInfo

# @BotFather'dan olgan tokeningizni shu yerga qo'ying
BOT_TOKEN = "8957524746:AAGMH7SWF9zaAuRu6RfPTHDaYdHQIQaRM_0"

# Telegram bergan Mini App havolangiz
WEB_APP_URL = "https://finance-tracker-five-neon.vercel.app"

# Botni ishga tushiramiz
bot = telebot.TeleBot(BOT_TOKEN)

# /start buyrug'i kelganda ishlaydigan qism
@bot.message_handler(commands=['start'])
def send_welcome(message):
    # Mini App-ni ochadigan tugma yaratamiz
    markup = InlineKeyboardMarkup()
    
    # Tugmaga nom beramiz va unga WebApp havolasini ulaymiz
    web_app_button = InlineKeyboardButton(
        text="🪙 Coin Tracker-ni ochish", 
        web_app=WebAppInfo(url=WEB_APP_URL)
    )
    
    # Tugmani shaklga qo'shamiz
    markup.add(web_app_button)
    
    # Foydalanuvchiga xabar va tugmani yuboramiz
    user_name = message.from_user.first_name
    bot.reply_to(
        message, 
        f"Salom {user_name}!\nCoin Tracker botiga xush kelibsiz. Ilovani ochish uchun quyidagi tugmani bosing 👇", 
        reply_markup=markup
    )

# Botni tinimsiz ishlab turishi uchun ishga tushiramiz
if __name__ == "__main__":
    print("Telebot muvaffaqiyatli ishga tushdi...")
    bot.infinity_polling()