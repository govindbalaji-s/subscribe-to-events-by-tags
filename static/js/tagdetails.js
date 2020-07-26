/**
 * Renders the TagDetails screen
 * props.data : data returned by API
 * props.user
 * props.onUnfollow
 * props.onFollow
 * props.onEventSubscribe
 * props.onEventUnsubscribe
 * props.onEventDetails
 */

function TagDetails(props) {
    let followUnfollowButton;
    if (props.data[APIIsFollowingField] == "true") {
        followUnfollowButton = (
            <Button onClick={(e) => {props.onUnfollow(props.data[APITagNameField]);}}>
                Unfollow
            </Button>
        );
    }
    else {
        followUnfollowButton = (
            <Button onClick={e => {props.onFollow(props.data[APITagNameField]);}}>
                Follow
            </Button>
        );
    }
    const combinedItems = props.data.taggedEventsData.map(event => ({
        item: event,
        otherProps: isSubscribed(props.user, event[APIEventIDField]) ? {
            actionLabel: "Unsubscribe",
            onAction: props.onEventUnsubscribe,
            onDetails: props.onEventDetails
        } : {
            actionLabel: "Subscribe",
            onAction: props.onEventSubscribe,
            onDetails: props.onEventDetails
        }
    }));
    console.log(followUnfollowButton);
    return (
        <React.Fragment>
            <div>
                Tag: {props.data[APITagNameField]}
            </div>
            <div>
                {followUnfollowButton}
            </div>
            <div>
                Followed by {props.data[APINoFollowersField]} people.
            </div>
            <div>
                <PagedList  header = {<h2>Events tagged:</h2>}
                                items = {combinedItems}
                                itemClass = {EventListRow}
                                keyFn = {eventKeyFn}
                />
            </div>
        </React.Fragment>
    )
}