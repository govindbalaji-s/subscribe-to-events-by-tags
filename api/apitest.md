user:
/user/get

create event:
$.post('/api/event/create', JSON.stringify({name: "hi", venue:"engo", time:"123456", duration:"12343", tags:["first"]}), d=>{console.log(d);}, 'json')

subscribe:
POST /api/event/subscribe/eventid