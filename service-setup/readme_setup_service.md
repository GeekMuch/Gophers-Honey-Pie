# Copy the service to correct location with command below: 

sudo cp gophers-start.service macchanger.service /lib/systemd/system/


# Copy the rc.local file to correct location with command below:

sudo cp rc.local /etc/

#####################################
# Create a backup of imange commands#
#####################################

# Make a backup of SD-card img

sudo dd if=/dev/mmcblk0 of=/raspberry-image.img

# After img creation, shrik the img with pishrink
# https://github.com/Drewsif/PiShrink

sudo pishrink.sh raspberry-image.img rpi-img-shrink.img


