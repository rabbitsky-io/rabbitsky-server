# Rabbit Sky Server Changelog

## Version 0.1.5
- Change position to INT instead of FLOAT, saving bandwidth.
- New Message Format, using regex, saving bandwidth.

## Version 0.1.4
- Add More Commands, /skyflash to change the color of the sky, flashing. /botadd to fill the server with bot that doing nothing, and /botremove to remove all the bots.
- Change players color from RGB to HSL instead.
- Change -password to -admin-password. We planned to use -password for password protected server like VIP only or something.

## Version 0.1.3
- Fix chat - when using comma, chat got cut

## Version 0.1.2
- Fix last player in init not sent when admin set sky

## Version 0.1.1
- Adding validation on Add Player Channel, so integer won't overflow.
- Adding admin command, for now only /admin for set the admin password, and /sky to change sky color using hexcode like #FFFFFF or rgb code like rgb(255,255,255)

## Version 0.1
- First Release