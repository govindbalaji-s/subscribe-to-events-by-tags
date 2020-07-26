/**
 * Returns whether user is following tagName
 */
function isFollowing(user, tagName) {
    const followingTags = user[APIUserTagsField];
    return followingTags.includes(tagName);
}

function isSubscribed(user, eventID) {
    const subscribedEventsIds = user[APIUserSubscribedEventsField];
    return subscribedEventsIds.includes(eventID);
}

/* Fetch full event details from the given id and call callbac
*/
function fetchEvent(eventId, callback) {
    $.get(`/api/event/get/${eventId}`, data => {
        if(data.result != 'success') {
            console.log(`Couldn't fetch event ${eventid}`);
        }
        else {
            callback(data.data);
        }
    }, 'json');
}

const tagKeyFn = tag => tag;
const eventKeyFn = event => event[APIEventIDField];