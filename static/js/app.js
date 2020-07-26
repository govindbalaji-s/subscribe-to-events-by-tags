'use strict';

const DASHBOARD_SCREEN = 0,
      TAG_DETAILS_SCREEN = 1;

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            screen: DASHBOARD_SCREEN,
            user: null
        }
        this.updateUser = this.updateUser.bind(this);
        this.onDashboardNav = this.onDashboardNav.bind(this);

        this.unfollowHandler = this.unfollowHandler.bind(this);
        this.onFollow = this.onFollow.bind(this);
        this.onTagDetails = this.onTagDetails.bind(this);

        this.onEventDetails = this.onEventDetails.bind(this);
        this.onEventSubscribe = this.onEventSubscribe.bind(this);
        this.onEventUnsubscribe = this.onEventUnsubscribe.bind(this);

        this.updateUser();
    }

    updateState() {
        this.updateUser();
        if(this.state.screen == DASHBOARD_SCREEN) {
            //do nothing
        }
        else if(this.state.screen == TAG_DETAILS_SCREEN) {
            this.onTagDetails(this.state.tagscreen[APITagNameField]);
        }
    }

    // Calls API to get the user object
    updateUser() {
        $.get('/user/get', resp => {
            if(resp.result == 'success') {
                resp.data.subscribedEventsData = [];
                this.setState({user: resp.data}, ()=>{
                    //subsrcibed events ids are only there, fetch all event details
                    resp.data[APIUserSubscribedEventsField].forEach((eventId, i) => {
                        fetchEvent(eventId, data => {
                            this.setState(prevState => {
                                prevState.user.subscribedEventsData[i] = data;
                                return prevState;
                            });
                        });
                    });
                });
            }
            else {
                this.setState({user:null});
            }
        });
    }

    onDashboardNav() {
        this.setState({screen: DASHBOARD_SCREEN}, this.updateState);
    }

    onEventDetails(eventid) {

    }

    onEventSubscribe(eventid) {
        $.post(`/api/event/subscribe/${eventid}`, resp => {
            console.log(`Subscribing ${eventid} was ${resp.result}`);
            this.updateState();
        }, 'json');
    }

    onEventUnsubscribe(eventid) {
        console.log(this);
        $.post(`/api/event/unsubscribe/${eventid}`, resp => {
            console.log(`Unsubscribing ${eventid} was ${resp.result}`);
            console.log(this);
            this.updateState();
        }, 'json');
    }
    // Calls API to unfollow a tag
    unfollowHandler(tagName) {
        console.log("Trying to unfollow"+tagName);
        $.post(`/api/tag/unfollow/${tagName}`, resp => {
            if(resp.result == 'success') {
                console.log(`Tag ${tagName} unfollowed.`);
            }
            this.updateState();
        }, 'json');
    }

    onFollow(tagName) {
        console.log("Trying to follow"+tagName);
        $.post(`/api/tag/follow/${tagName}`, resp => {
            if(resp.result == 'success') {
                console.log(`Tag ${tagName} followed.`);
            }
            this.updateState();
        }, 'json');
    }

    onTagDetails(tagName) {
        console.log("Getting tag"+tagName);
        $.get(`/api/tag/get/${tagName}`, resp => {
            if(resp.result == 'success') {
                resp.data.taggedEventsData = [];
                this.setState({
                    screen: TAG_DETAILS_SCREEN,
                    tagscreen: resp.data
                }, () => {
                    resp.data[APITaggedEventsField].forEach((eventId, i) => {
                        fetchEvent(eventId, data => {
                            this.setState(prevState => {
                                prevState.tagscreen.taggedEventsData[i] = data;
                                return prevState;
                            });
                        });
                    });
                });
            }
        });
    }

    render() {
        let screenComponent;
        if(this.state.screen == DASHBOARD_SCREEN) {
            screenComponent =  (
                <Dashboard
                    user = {this.state.user}
                    updateUser = {this.updateUser}
                    unfollowHandler = {this.unfollowHandler}
                    onTagDetails = {this.onTagDetails}
                    onEventDetails = {this.onEventDetails}
                    onEventSubscribe = {this.onEventSubscribe}
                    onEventUnsubscribe = {this.onEventUnsubscribe}
                />
            );
        }
        else if(this.state.screen == TAG_DETAILS_SCREEN) {
            screenComponent = (
                <TagDetails
                    user = {this.state.user}
                    data = {this.state.tagscreen}
                    onUnfollow = {this.unfollowHandler}
                    onFollow = {this.onFollow}
                    onEventSubscribe = {this.onEventSubscribe}
                    onEventUnsubscribe = {this.onEventUnsubscribe}
                    onEventDetails = {this.onEventDetails}
                />
            );
        }
        return (
            <React.Fragment>
                <NavBar
                    user = {this.state.user}
                    onDashboardNav = {this.onDashboardNav}
                />
                {this.state.user ? screenComponent : null}
            </React.Fragment>
        )
    }
}
ReactDOM.render(<App/>, document.getElementById('root'));