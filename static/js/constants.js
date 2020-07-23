const   APIUserEmailField            = "email",
        APITagNameField              = "name",
        APIUserTagsField             = "tags",
        APIEventsField               = "events",
        APIFollowersField            = "followers",

        APIEventIDField              = "eventID",
        APIEventNameField            = "name",
        APIEventVenueField           = "venue",
        APIEventTimeField            = "time",
        APIEventDurationField        = "duration",
        APIEventTagsField            = "tags",
        APIEventNoSubscribersField   = "noOfSubscribers",
        APIEventIsSubscribedField    = "isSubscribed",

        APIEventCreatorField         = "creator",
        APIUserCreatedEventsField    = "createdEvents",
        APIUserSubscribedEventsField = "subscribedEvents",

        TimeFormat = "02-01-2006 15:04 (IST)";

function timestampToString(timestamp) {
    let x = new Date(timestamp * 1000);
    let dd = x.getDate();
    dd < 10 && (dd = '0' + dd);
    let mm = x.getMonth() + 1;
    mm < 10 && (mm = '0' + mm);
    let yyyy = x.getFullYear();
    let hh = x.getHours();
    if (hh == 0)
        hh = '00';
    else if(hh < 10)
        hh = '0' + hh;
    let min = x.getMinutes();
    if (min == 0)
        min = '00';
    else if(min < 10)
        min = '0' + min;

    return `${dd}/${mm}/${yyyy} ${hh}:${min}`;    
}