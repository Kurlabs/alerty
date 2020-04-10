### How works?
You should send a pub/sub message to messenger topic with the next structure:

* Data:
No needed right now

* Attributes:
number: phone number like +18017871055
`Sms` (_bool_): Flag to indicate if SMS notification has to be made
`Email` (_bool_): Flag to indicate if Email notification has to be made
`ObjectType`(_string_): Object type. ie: _domain_ or _socket_
`Object`(_string_): Object value. ie: https://alerty.online, http://12.34.56.78:1234
`Event`(_string_): Event occurred. ie: _uptime_, _downtime_
`EmailAddresses`(_string_): Email addresses to be notified. This value must be comma separated. ie: "_me@domain.com_,_other@domain.com_"
`PhoneNumbers`(_string_) Phone numbers to be notified. This value must be comma separated. ie: "_+18017871055_,_+573195206895_"
`UserInfo`(_string_): User info to be include in the messages templates