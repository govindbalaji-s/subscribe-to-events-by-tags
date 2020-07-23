/*
** To browse through list of following tags and unfollow them
** Required props:
**      props.tags -> list of tags
**      props.unFollowHandler -> function(TagNameField) to call for unfollow event
**      props.onTagDetails -> move to tag page
*/
function FollowingTags(props) {
    return (
        <PagedList  header = {<h2>Tags you follow:</h2>}
                    items = {props.tags}
                    itemClass = {TagListRow}
                    keyFn = {tag => tag}
                    otherProps = {{
                                    actionLabel: "Unfollow",
                                    onUnfollow: props.unfollowHandler,
                                    onTagDetails: props.onTagDetails
                                }}
        />
    )
}

/*
 ** List of subscribed events
 ** Required props:
 **     props.eventIDs
 **     props.onEventDetails
 **     props.onEventToggleSubscription
 */
function SubscribedEvents(props) {
    console.log(props.events);
    return (
        <PagedList  header = {<h2>Events subscribed:</h2>}
                    items = {props.events}
                    itemClass = {EventListRow}
                    keyFn = {event => event[APIEventIDField]}
                    otherProps = {{
                        onDetails: props.onEventDetails, // TODO: move to events page
                        onAction: props.onEventToggleSubscription,
                        actionLabel: "Unsubscribe"
                    }}
        />
    )
}

/*
** The entire dashboard
** Required props:
**      props.user -> user object
**      props.unfollowHandler -> function(TagNameField) to call when unfollowing a tag
**      props.onFetchEvents -> function([EventID]) to call to fetch events
**      props.onEventDetails -> function(eventid) to call to move to event page
**      props.onEventSubscribe,
**      props.onEventUnsubscribe -> function(eventid) to call to toggle subscription to the event
*/
class Dashboard extends React.Component {
    render() {
        if(!this.props.user) { //check if no user signed in
            return <LoginButton />;
        }
        else {
            console.log(this.props.user);
            return (
                <div>
                    Hello {this.props.user[APIUserEmailField]}!
                    <br/>
                    <LogoutButton />
                    <FollowingTags
                        tags={this.props.user[APIUserTagsField]}
                        unfollowHandler = {this.props.unfollowHandler}
                        onTagDetails = {this.props.onTagDetails}
                    />
                    <SubscribedEvents
                        eventIDs = {this.props.user[APIUserSubscribedEventsField]}
                        events = {this.props.user.subscribedEventsData}
                        onFetchEvents = {this.props.onFetchEvents}
                        onEventDetails = {this.props.onEventDetails}
                        onEventToggleSubscription = {this.props.onEventUnsubscribe}
                    />
                </div>
                );
        }
    }
}