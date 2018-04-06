from lcdproc.server import Server
import requests
import time
import unicodedata

#===================================================
class Display():
    def __init__(self):
        self.lcd = Server("localhost", debug=False)
        self.lcd.start_session()
        self.screen = self.lcd.add_screen("Volumio")
        self.screen.set_width(16)
        self.screen.set_height(2)
        self.screen.set_heartbeat("off")
        self.screen.set_backlight("off")
        self.artist_widget = self.screen.add_string_widget("ArtistWidget", text="<Artist>")
        self.title_widget = self.screen.add_scroller_widget("TitleWidget", text="<Title>", speed=4, top=2, right=16)
        
        self.stop = True

    def on_play(self):
        if self.stop:
            self.screen.set_backlight("on")
            self.stop = False 

    def on_stop(self):
        if not self.stop:
            self.screen.set_backlight("off")
            self.stop = True

    def artist(self, text):
        self.artist_widget.set_text(text)

    def title(self, text):
        self.title_widget.set_text(text)
        

#===================================================
class Player():
    def __init__(self):
        self.url = "http://localhost:3000"
        self.status = "stop"
        self.title = ""
        self.artist = ""
        self.player_data = {}
        self.dirty = False

    def update(self):
        command = "/api/v1/getstate"
        response = requests.get(self.url + command)
        data = response.json()
        self.dirty = data != self.player_data
        if not self.dirty:
            return False

        self.player_data = data
        self.status = unicodedata.normalize('NFKD', data.get("status", u"stop")).encode('ascii', 'ignore')
        self.title = unicodedata.normalize('NFKD', data.get("title", u"")).encode('ascii', 'ignore')
        self.artist = unicodedata.normalize('NFKD', data.get("artist", u"")).encode('ascii', 'ignore')

        return True

#===================================================
def main():
    display = Display()
    player = Player()
    
    while True:
        if not player.update():
            time.sleep(0.5)
            continue

        display.artist(player.artist)
        display.title(player.title)

        if player.status == "play":
            display.on_play()
        else:
            display.on_stop()
        
        time.sleep(0.5)

if __name__ == "__main__":
    main()
