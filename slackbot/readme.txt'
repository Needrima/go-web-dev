Instructions to get tokens:

create a workspace on slack
visit api.slack.com and access your workspace
click create new app (choose from scratch and specify the name of the app/bot and the workspace you want it to be)
create the app
navigate to 'Socket mode' and 'Enable socket mode' and generate socket token which is also app token
copy the token and save somewhere secret maybe as an environmental variable
navigate to'Event subscription tab' and 'Enable events'
navigate to 'oAuth and permissions' tab
Add oAuth scopes:
    app_mentions:read
    channels:history
    channels:read
    chat:write
    im:history
    im:read
    im:write
    mpim:history
    mpim:read
    mpim:write
go back to 'events subscription' tab and 're-enable events' if disabled
subsribe to the following bot events:
    app_mention
    message.im
    message.mpim
    message.groups
    message.channels
navigate back to 'oAuth and permissions' tab 
click the 'Install to workspace' button
copy the 'bot token' and store somwhere safe like an environmental variable
NOTE: if any change is made after installing, you have to 'reinstall to workspace'