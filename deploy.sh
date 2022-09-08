#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o meddle .
rsync -auv -e "ssh -i /Users/mac/downloads/meddle.pem" ./server/templates/* ubuntu@ec2-3-143-47-141.us-east-2.compute.amazonaws.com:/home/ubuntu/server/templates/
rsync -auv -e "ssh -i /Users/mac/downloads/meddle.pem" ./meddle ubuntu@ec2-3-143-47-141.us-east-2.compute.amazonaws.com:/home/ubuntu/
ssh -i /Users/mac/downloads/meddle.pem ubuntu@ec2-3-143-47-141.us-east-2.compute.amazonaws.com "sudo systemctl restart meddle.service"
sleep 5
ssh -i /Users/mac/downloads/meddle.pem ubuntu@ec2-3-143-47-141.us-east-2.compute.amazonaws.com "sudo systemctl status meddle.service"